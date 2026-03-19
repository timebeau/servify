import EventEmitter from 'eventemitter3';
import type { AuthProvider } from './contracts/auth-provider';
import { ServifyError } from './contracts/errors';
import type { Transport, TransportConnectOptions, TransportSendOptions, ReconnectPolicy, TransportState } from './contracts/transport';
import { WSMessage, ServifyEventMap } from './types';

export interface WebSocketManagerOptions {
  url: string;
  protocols?: string | string[];
  reconnectAttempts?: number;
  reconnectDelay?: number;
  heartbeatInterval?: number;
  debug?: boolean;
  reconnectPolicy?: ReconnectPolicy;
  authProvider?: AuthProvider;
  onTokenRefreshRequired?: () => Promise<void>;
}

type NormalizedWebSocketManagerOptions = Omit<
  Required<WebSocketManagerOptions>,
  'authProvider'
> & {
  authProvider?: AuthProvider;
};

export class WebSocketManager extends EventEmitter<ServifyEventMap> implements Transport<WSMessage, WSMessage> {
  private ws: WebSocket | null = null;
  private options: NormalizedWebSocketManagerOptions;
  private reconnectAttempts = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private heartbeatTimer: NodeJS.Timeout | null = null;
  private isManualClose = false;
  private subscribers = new Set<(message: WSMessage) => void>();
  readonly kind = 'websocket';
  state: TransportState = 'idle';

  constructor(options: WebSocketManagerOptions) {
    super();

    this.options = {
      protocols: [],
      reconnectAttempts: 5,
      reconnectDelay: 1000,
      heartbeatInterval: 30000,
      debug: false,
      reconnectPolicy: {
        maxAttempts: options.reconnectAttempts ?? 5,
        baseDelayMs: options.reconnectDelay ?? 1000,
        backoffFactor: 2,
        maxDelayMs: 30000,
      },
      authProvider: options.authProvider,
      onTokenRefreshRequired: options.onTokenRefreshRequired ?? (async () => undefined),
      ...options
    };
  }

  connect(_options?: TransportConnectOptions): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.state = 'connected';
        resolve();
        return;
      }

      this.isManualClose = false;
      this.state = 'connecting';
      this.log('正在连接 WebSocket...', this.options.url);

      try {
        this.ws = new WebSocket(this.options.url, this.options.protocols);
      } catch (error) {
        this.state = 'error';
        reject(error);
        return;
      }

      this.ws.onopen = () => {
        this.log('WebSocket 连接成功');
        this.reconnectAttempts = 0;
        this.state = 'connected';
        this.startHeartbeat();
        this.emit('connected');
        resolve();
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          this.log('解析消息失败:', error);
          this.emit('error', new Error('Invalid message format'));
        }
      };

      this.ws.onclose = (event) => {
        this.log('WebSocket 连接关闭:', event.code, event.reason);
        this.stopHeartbeat();
        this.state = this.isManualClose ? 'closed' : 'idle';
        this.emit('disconnected', event.reason || '连接关闭');

        if (!this.isManualClose && this.reconnectAttempts < this.options.reconnectPolicy.maxAttempts) {
          this.state = 'reconnecting';
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = (event) => {
        this.log('WebSocket 错误:', event);
        this.state = 'error';
        const err = new ServifyError('WebSocket connection error', {
          code: 'transport_unavailable',
          retryable: true,
          details: { url: this.options.url },
        });
        this.emit('error', err);
        reject(err);
      };
    });
  }

  async disconnect(): Promise<void> {
    this.isManualClose = true;
    this.state = 'closed';
    this.stopHeartbeat();

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  async send(message: WSMessage, _options?: TransportSendOptions): Promise<void> {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.log('WebSocket 未连接，无法发送消息');
      const err = new ServifyError('WebSocket not connected', {
        code: 'transport_disconnected',
        retryable: true,
      });
      this.emit('error', err);
      throw err;
    }

    try {
      this.ws.send(JSON.stringify(message));
      this.log('发送消息:', message);
    } catch (error) {
      this.log('发送消息失败:', error);
      const err = new ServifyError('Failed to send message', {
        code: 'transport_unavailable',
        cause: error,
        retryable: true,
      });
      this.emit('error', err);
      throw err;
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  subscribe(handler: (message: WSMessage) => void): () => void {
    this.subscribers.add(handler);
    return () => {
      this.subscribers.delete(handler);
    };
  }

  private handleMessage(message: WSMessage): void {
    this.log('收到消息:', message);
    for (const subscriber of this.subscribers) {
      subscriber(message);
    }

    switch (message.type) {
      case 'message':
        this.emit('message', message.data);
        break;
      case 'session_update':
        this.emit('session_updated', message.data);
        break;
      case 'agent_status':
        if (message.data.type === 'assigned') {
          this.emit('agent_assigned', message.data.agent);
        } else if (message.data.type === 'typing') {
          this.emit('agent_typing', message.data.typing);
        }
        break;
      case 'error':
        this.emit('error', new Error(message.data.message || 'Unknown error'));
        break;
      case 'system':
        // 处理系统消息，如心跳响应
        if (message.data?.type === 'pong') {
          // 心跳响应处理
          this.log('收到心跳响应');
        }
        break;
      default:
        this.log('未知消息类型:', message.type);
    }
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    this.emit('reconnecting', this.reconnectAttempts);

    const factor = this.options.reconnectPolicy.backoffFactor ?? 2;
    const maxDelay = this.options.reconnectPolicy.maxDelayMs ?? 30000;
    const delay = Math.min(
      this.options.reconnectPolicy.baseDelayMs * Math.pow(factor, this.reconnectAttempts - 1),
      maxDelay
    );
    this.log(`${delay}ms 后重连 (第 ${this.reconnectAttempts}/${this.options.reconnectPolicy.maxAttempts} 次)`);

    this.reconnectTimer = setTimeout(() => {
      this.connect().catch(() => {
        // 重连失败，继续尝试或放弃
      });
    }, delay);
  }

  private startHeartbeat(): void {
    this.stopHeartbeat();

    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected()) {
        void this.send({
          type: 'system',
          data: { type: 'ping', timestamp: new Date().toISOString() }
        }).catch(() => undefined);
      }
    }, this.options.heartbeatInterval);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private log(...args: any[]): void {
    if (this.options.debug) {
      console.log('[ServifyWS]', ...args);
    }
  }
}
