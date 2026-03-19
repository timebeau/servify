import type { AuthProvider } from './auth-provider';

export type TransportState =
  | 'idle'
  | 'connecting'
  | 'connected'
  | 'reconnecting'
  | 'closed'
  | 'error';

export interface ReconnectPolicy {
  maxAttempts: number;
  baseDelayMs: number;
  maxDelayMs?: number;
  backoffFactor?: number;
}

export interface TransportEnvelope<TPayload = unknown> {
  type: string;
  data: TPayload;
  timestamp?: string;
}

export interface TransportConnectOptions {
  signal?: AbortSignal;
}

export interface TransportSendOptions {
  requireAck?: boolean;
}

export interface TransportOptions {
  debug?: boolean;
  reconnectPolicy?: ReconnectPolicy;
  authProvider?: AuthProvider;
  onTokenRefreshRequired?: () => Promise<void>;
}

export interface Transport<TInbound = unknown, TOutbound = unknown> {
  readonly kind: string;
  readonly state: TransportState;
  connect(options?: TransportConnectOptions): Promise<void>;
  disconnect(): Promise<void> | void;
  send(message: TOutbound, options?: TransportSendOptions): Promise<void> | void;
  isConnected(): boolean;
  subscribe(handler: (message: TInbound) => void): () => void;
}
