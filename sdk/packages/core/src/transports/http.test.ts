import { describe, expect, it, vi } from 'vitest';

import { HttpTransport } from './http';

describe('HttpTransport', () => {
  it('tracks connection state and allows subscriptions', async () => {
    const transport = new HttpTransport<string, string>();
    const handler = vi.fn();

    const unsubscribe = transport.subscribe(handler);
    expect(transport.isConnected()).toBe(false);

    await transport.connect();
    expect(transport.isConnected()).toBe(true);
    expect(transport.state).toBe('connected');

    await transport.send('ping');
    unsubscribe();

    await transport.disconnect();
    expect(transport.state).toBe('closed');
  });
});

