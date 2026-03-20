export interface RetryBackoffPolicy {
  maxAttempts: number;
  baseDelayMs: number;
  maxDelayMs?: number;
  backoffFactor?: number;
  retryableStatusCodes?: number[];
}

export interface RetryDecisionContext {
  attempt: number;
  statusCode?: number;
  isNetworkError?: boolean;
}

const DEFAULT_RETRY_POLICY: RetryBackoffPolicy = {
  maxAttempts: 3,
  baseDelayMs: 250,
  maxDelayMs: 5000,
  backoffFactor: 2,
  retryableStatusCodes: [408, 409, 425, 429, 500, 502, 503, 504],
};

export function normalizeRetryBackoffPolicy(
  policy?: Partial<RetryBackoffPolicy>,
): RetryBackoffPolicy {
  return {
    maxAttempts: policy?.maxAttempts ?? DEFAULT_RETRY_POLICY.maxAttempts,
    baseDelayMs: policy?.baseDelayMs ?? DEFAULT_RETRY_POLICY.baseDelayMs,
    maxDelayMs: policy?.maxDelayMs ?? DEFAULT_RETRY_POLICY.maxDelayMs,
    backoffFactor: policy?.backoffFactor ?? DEFAULT_RETRY_POLICY.backoffFactor,
    retryableStatusCodes: policy?.retryableStatusCodes ?? DEFAULT_RETRY_POLICY.retryableStatusCodes,
  };
}

export function shouldRetryRequest(
  context: RetryDecisionContext,
  policy: RetryBackoffPolicy,
): boolean {
  if (context.attempt >= policy.maxAttempts) {
    return false;
  }

  if (context.isNetworkError) {
    return true;
  }

  return context.statusCode !== undefined && policy.retryableStatusCodes?.includes(context.statusCode) === true;
}

export function computeRetryDelay(policy: RetryBackoffPolicy, attempt: number): number {
  const exponent = Math.max(attempt - 1, 0);
  const delay = policy.baseDelayMs * Math.pow(policy.backoffFactor ?? 2, exponent);
  return Math.min(delay, policy.maxDelayMs ?? delay);
}
