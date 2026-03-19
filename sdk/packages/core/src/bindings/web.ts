import { StaticCapabilitySet } from '../contracts/capability';

export function createWebCapabilitySet(): StaticCapabilitySet {
  return new StaticCapabilitySet([
    { name: 'chat', enabled: true, version: '1' },
    { name: 'realtime', enabled: true, version: '1' },
    { name: 'knowledge', enabled: true, version: '1' },
    { name: 'remote_assist', enabled: false, version: 'reserved' },
    { name: 'voice', enabled: false, version: 'reserved' },
  ]);
}
