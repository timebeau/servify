import { describe, expect, it } from 'vitest';

import { StaticCapabilitySet, negotiateCapabilities } from './capability';

describe('StaticCapabilitySet', () => {
  it('returns copies and enabled checks', () => {
    const set = new StaticCapabilitySet([
      { name: 'chat', enabled: true, version: '1' },
      { name: 'voice', enabled: false, version: 'reserved' },
    ]);

    expect(set.has('chat')).toBe(true);
    expect(set.has('voice')).toBe(false);
    expect(set.get('voice')?.version).toBe('reserved');

    const entries = set.all();
    entries.push({ name: 'knowledge', enabled: true, version: '1' });

    expect(set.all()).toHaveLength(2);
  });

  it('negotiates enabled, reserved, and unsupported capabilities', () => {
    const set = new StaticCapabilitySet([
      { name: 'chat', enabled: true, version: '1' },
      { name: 'voice', enabled: false, version: 'reserved' },
    ]);

    const negotiated = set.negotiate([
      { name: 'chat' },
      { name: 'voice' },
      { name: 'remote_assist' },
    ]);

    expect(negotiated.granted).toEqual([
      { name: 'chat', enabled: true, version: '1' },
    ]);
    expect(negotiated.rejected).toEqual([
      {
        request: { name: 'voice' },
        reason: 'disabled',
        descriptor: { name: 'voice', enabled: false, version: 'reserved' },
      },
      {
        request: { name: 'remote_assist' },
        reason: 'unsupported',
      },
    ]);
  });
});

describe('negotiateCapabilities', () => {
  it('keeps returned descriptors isolated from source arrays', () => {
    const source = [{ name: 'knowledge', enabled: true, version: '1' }] as const;
    const result = negotiateCapabilities([...source], [{ name: 'knowledge' }]);

    result.granted[0].metadata = { changed: true };

    expect(source[0]).toEqual({ name: 'knowledge', enabled: true, version: '1' });
  });
});
