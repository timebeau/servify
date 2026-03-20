import type { BearerTokenAuthProvider } from './contracts/auth';
import type { IdempotencyKeyProvider } from './contracts/request-pipeline';
import { normalizeRetryBackoffPolicy } from './contracts/retry';

export interface ApiClientAutomationExample {
  name: string;
  description: string;
  auth: BearerTokenAuthProvider | { kind: 'api_key'; keyName: string };
  retryPolicy: ReturnType<typeof normalizeRetryBackoffPolicy>;
  idempotency: IdempotencyKeyProvider;
}

export const automationExamples: ApiClientAutomationExample[] = [
  {
    name: 'admin-ticket-export',
    description: '后台定时导出工单，需要稳定重试与幂等键以避免重复任务。',
    auth: {
      kind: 'api_key',
      keyName: 'x-servify-admin-key',
    },
    retryPolicy: normalizeRetryBackoffPolicy({
      maxAttempts: 5,
      baseDelayMs: 500,
    }),
    idempotency: {
      generate: ({ path }) => `admin-export:${path}`,
    },
  },
  {
    name: 'bot-session-sync',
    description: '机器人侧同步会话状态，需要 bearer token 与短延迟退避。',
    auth: {
      kind: 'bearer',
      getAccessToken: async () => ({ token: 'reserved', tokenType: 'Bearer' }),
      getHeaders: async () => ({ authorization: 'Bearer reserved' }),
    },
    retryPolicy: normalizeRetryBackoffPolicy({
      maxAttempts: 4,
      baseDelayMs: 200,
    }),
    idempotency: {
      generate: ({ path }) => `bot-sync:${path}`,
    },
  },
];
