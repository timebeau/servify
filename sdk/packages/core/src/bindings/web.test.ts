import { describe, expect, it } from 'vitest';

import { createWebServifySDK } from './web-sdk';
import { createWebCapabilitySet } from './web';

describe('web bindings', () => {
  it('exposes the default web capability set', () => {
    const capabilities = createWebCapabilitySet();

    expect(capabilities.has('chat')).toBe(true);
    expect(capabilities.has('realtime')).toBe(true);
    expect(capabilities.has('remote_assist')).toBe(true);
    expect(capabilities.get('voice')?.enabled).toBe(false);
  });

  it('creates a web sdk client with remote assistance capability attached', () => {
    const sdk = createWebServifySDK({
      apiUrl: 'http://localhost:8080',
      autoConnect: false,
    });

    expect(sdk.capabilities.has('chat')).toBe(true);
    expect(sdk.capabilities.get('remote_assist')?.enabled).toBe(true);
    expect(sdk.capabilities.get('remote_assist')?.version).toBe('1');
  });
});
