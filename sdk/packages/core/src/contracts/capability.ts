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

export interface CapabilityRequest {
  name: CapabilityName;
}

export interface CapabilityNegotiationRejection {
  request: CapabilityRequest;
  reason: 'disabled' | 'unsupported';
  descriptor?: CapabilityDescriptor;
}

export interface CapabilityNegotiationResult {
  granted: CapabilityDescriptor[];
  rejected: CapabilityNegotiationRejection[];
}

export interface CapabilitySet {
  all(): CapabilityDescriptor[];
  has(name: CapabilityName): boolean;
  get(name: CapabilityName): CapabilityDescriptor | undefined;
  negotiate(requested: CapabilityRequest[]): CapabilityNegotiationResult;
}

export function negotiateCapabilities(
  available: CapabilityDescriptor[],
  requested: CapabilityRequest[],
): CapabilityNegotiationResult {
  const granted: CapabilityDescriptor[] = [];
  const rejected: CapabilityNegotiationRejection[] = [];

  for (const request of requested) {
    const descriptor = available.find((entry) => entry.name === request.name);
    if (!descriptor) {
      rejected.push({ request, reason: 'unsupported' });
      continue;
    }

    if (!descriptor.enabled) {
      rejected.push({ request, reason: 'disabled', descriptor });
      continue;
    }

    granted.push({ ...descriptor });
  }

  return { granted, rejected };
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

  negotiate(requested: CapabilityRequest[]): CapabilityNegotiationResult {
    return negotiateCapabilities(this.entries, requested);
  }
}
