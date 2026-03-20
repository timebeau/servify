import { afterEach, describe, expect, it, vi } from 'vitest';

import { WebSocketManager } from './websocket';

class FakeWebSocket {
  static instances: FakeWebSocket[] = [];
  static OPEN = 1;

  readonly url: string;
  readonly protocols?: string | string[];
  readyState = FakeWebSocket.OPEN;
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onclose: ((event: { code: number; reason: string }) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;

  constructor(url: string, protocols?: string | string[]) {
    this.url = url;
    this.protocols = protocols;
    FakeWebSocket.instances.push(this);
  }

  send(): void {}

  close(): void {
    this.onclose?.({ code: 1000, reason: 'closed' });
  }

  open(): void {
    this.onopen?.();
  }
}

describe('WebSocketManager', () => {
  afterEach(() => {
    FakeWebSocket.instances = [];
    vi.unstubAllGlobals();
  });

  it('adds access token from auth provider before connecting', async () => {
    vi.stubGlobal('WebSocket', FakeWebSocket);

    const manager = new WebSocketManager({
      url: 'ws://localhost:8080/ws?customer_id=1',
      authProvider: {
        getToken: async () => ({ accessToken: 'token-123' }),
      },
    });

    const connectPromise = manager.connect();
    await vi.waitFor(() => expect(FakeWebSocket.instances).toHaveLength(1));
    FakeWebSocket.instances[0].open();
    await connectPromise;

    expect(FakeWebSocket.instances[0].url).toContain('access_token=token-123');
    expect(FakeWebSocket.instances[0].url).toContain('customer_id=1');
  });

  it('refreshes token when the current token is unavailable', async () => {
    vi.stubGlobal('WebSocket', FakeWebSocket);

    const onTokenRefreshRequired = vi.fn(async () => undefined);
    const refreshToken = vi.fn(async () => ({ accessToken: 'refreshed-token' }));

    const manager = new WebSocketManager({
      url: 'ws://localhost:8080/ws',
      authProvider: {
        getToken: async () => null,
        refreshToken,
      },
      onTokenRefreshRequired,
    });

    const connectPromise = manager.connect();
    await vi.waitFor(() => expect(FakeWebSocket.instances).toHaveLength(1));
    FakeWebSocket.instances[0].open();
    await connectPromise;

    expect(onTokenRefreshRequired).toHaveBeenCalledTimes(1);
    expect(refreshToken).toHaveBeenCalledTimes(1);
    expect(FakeWebSocket.instances[0].url).toContain('access_token=refreshed-token');
  });

  it('fails with auth_refresh_required when refresh hook cannot produce a token', async () => {
    vi.stubGlobal('WebSocket', FakeWebSocket);

    const manager = new WebSocketManager({
      url: 'ws://localhost:8080/ws',
      authProvider: {
        getToken: async () => null,
        refreshToken: async () => null,
      },
      onTokenRefreshRequired: async () => undefined,
    });

    await expect(manager.connect()).rejects.toMatchObject({
      code: 'auth_refresh_required',
    });
  });
});
