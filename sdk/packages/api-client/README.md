# @servify/api-client

保留包，用于未来服务端与服务端之间的 SDK 访问层。

预留场景：

- server-to-server API calls
- admin automation clients
- bot integrations

当前不实现具体 transport/client 逻辑，但已经固定以下 contract：

- server-side auth provider contract
- retry/backoff policy contract
- idempotency key 与 request middleware contract
- bot/admin automation 示例

## Exported Contracts

- `ServerAuthProvider` / `BearerTokenAuthProvider` / `ApiKeyAuthProvider`
- `RetryBackoffPolicy` / `normalizeRetryBackoffPolicy` / `computeRetryDelay`
- `ApiRequestMiddleware` / `IdempotencyKeyProvider`
- `automationExamples`

这些 contract 仅用于未来 server-side SDK 设计收口，当前不承诺具体 runtime 行为。
