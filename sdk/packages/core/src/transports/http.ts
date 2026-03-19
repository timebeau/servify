import type { Transport, TransportConnectOptions, TransportSendOptions, TransportState } from '../contracts/transport';

// HTTP transport is reserved for future request/response and long-polling scenarios.
export class HttpTransport<TInbound = unknown, TOutbound = unknown>
  implements Transport<TInbound, TOutbound> {
  readonly kind = 'http';
  state: TransportState = 'idle';

  async connect(_options?: TransportConnectOptions): Promise<void> {
    this.state = 'connected';
  }

  async disconnect(): Promise<void> {
    this.state = 'closed';
  }

  async send(_message: TOutbound, _options?: TransportSendOptions): Promise<void> {
    return;
  }

  isConnected(): boolean {
    return this.state === 'connected';
  }

  subscribe(_handler: (message: TInbound) => void): () => void {
    return () => undefined;
  }
}
