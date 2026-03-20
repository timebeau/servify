import type { ReconnectPolicy } from './transport';

export interface LegacyReconnectOptions {
  reconnectAttempts?: number;
  reconnectDelay?: number;
}

export interface ReconnectDecision {
  attempt: number;
  isManualClose: boolean;
}

const DEFAULT_RECONNECT_POLICY: ReconnectPolicy = {
  maxAttempts: 5,
  baseDelayMs: 1000,
  backoffFactor: 2,
  maxDelayMs: 30000,
};

export function normalizeReconnectPolicy(
  policy?: Partial<ReconnectPolicy>,
  legacy?: LegacyReconnectOptions,
): ReconnectPolicy {
  return {
    maxAttempts: policy?.maxAttempts ?? legacy?.reconnectAttempts ?? DEFAULT_RECONNECT_POLICY.maxAttempts,
    baseDelayMs: policy?.baseDelayMs ?? legacy?.reconnectDelay ?? DEFAULT_RECONNECT_POLICY.baseDelayMs,
    backoffFactor: policy?.backoffFactor ?? DEFAULT_RECONNECT_POLICY.backoffFactor,
    maxDelayMs: policy?.maxDelayMs ?? DEFAULT_RECONNECT_POLICY.maxDelayMs,
  };
}

export function computeReconnectDelay(policy: ReconnectPolicy, attempt: number): number {
  const exponent = Math.max(attempt - 1, 0);
  const delay = policy.baseDelayMs * Math.pow(policy.backoffFactor ?? 2, exponent);
  return Math.min(delay, policy.maxDelayMs ?? delay);
}

export function shouldReconnect(decision: ReconnectDecision, policy: ReconnectPolicy): boolean {
  if (decision.isManualClose) {
    return false;
  }

  return decision.attempt < policy.maxAttempts;
}
