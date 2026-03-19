export type CapabilityName =
  | 'chat'
  | 'realtime'
  | 'knowledge'
  | 'remote_assist'
  | 'voice';

export interface CapabilityDescriptor {
  name: CapabilityName;
  enabled: boolean;
  version?: string;
  metadata?: Record<string, unknown>;
}

export interface CapabilitySet {
  all(): CapabilityDescriptor[];
  has(name: CapabilityName): boolean;
  get(name: CapabilityName): CapabilityDescriptor | undefined;
}

export class StaticCapabilitySet implements CapabilitySet {
  private readonly entries: CapabilityDescriptor[];

  constructor(entries: CapabilityDescriptor[]) {
    this.entries = entries.slice();
  }

  all(): CapabilityDescriptor[] {
    return this.entries.slice();
  }

  has(name: CapabilityName): boolean {
    return this.entries.some((entry) => entry.name === name && entry.enabled);
  }

  get(name: CapabilityName): CapabilityDescriptor | undefined {
    return this.entries.find((entry) => entry.name === name);
  }
}
