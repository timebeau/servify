import { describe, expect, it } from 'vitest';

import { StaticCapabilitySet } from './capability';

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
});

