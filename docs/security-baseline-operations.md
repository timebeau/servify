# Security Baseline For Operations

本文件用于推进 `11 / T5 security-baseline-for-operations`，先把当前高风险接口、已有控制和最小上线基线收口成统一清单。

## 高风险接口盘点

### Management Surface 写接口

范围：

- `/api/*` 下的工单、客户、坐席、转接、集成、自动化、知识库等写操作
- `/api/v1/ai/knowledge/upload`
- `/api/v1/ai/knowledge/sync`
- `/api/v1/ai/weknora/enable`
- `/api/v1/ai/weknora/disable`
- `/api/v1/ai/circuit-breaker/reset`

风险：

- 会改变业务主数据、路由行为或 AI/knowledge 运行状态
- 错误授权会直接造成跨租户写入、配置污染或服务降级

当前控制：

- `AuthMiddleware`
- `EnforceRequestScope`
- `RequirePrincipalKinds("agent", "admin", "service")`
- `RequireResourcePermission(...)`
- 管理面写请求统一经过 `AuditMiddleware`

### Service Surface

范围：

- `/api/v1/metrics/ingest`

风险：

- 属于机器入口，若 token 泄漏可能被用于伪造监控数据或造成放大写入

当前控制：

- `AuthMiddleware`
- `EnforceRequestScope`
- `RequirePrincipalKinds("service")`
- 全局 HTTP rate limiting middleware 可按路径单独限流

### Public Surface

范围：

- `/public/portal/config`
- `/public/kb/*`
- `/public/csat/*`
- `/api/v1/ws`

风险：

- 匿名暴露，容易成为枚举、抓取、滥用或流量放大入口

当前控制：

- 路由表面已与 management / service 分离，避免复用后台 JWT 语义
- 全局 HTTP rate limiting middleware 已挂到基础路由

## 已落地的安全基线

### 认证与授权

- 管理面与 service 面统一走 `internal/platform/auth`
- `end_user` token 被显式拒绝进入 management surface
- request scope 会校验 token scope 与 header/query scope 冲突，阻止请求扩大作用域

### 审计

- 管理面写请求默认记录 actor、principal、tenant、workspace、resource、request metadata
- 审计查询接口：`GET /api/audit/logs`
- 审计中间件现已对 `password`、`secret`、`api_key`、`token` 等敏感字段做统一脱敏，避免把明文 secrets 写入审计库

### 速率限制

- `registerBaseMiddleware` 已统一挂载 `RateLimitMiddlewareFromConfig`
- 支持全局限流、按路径前缀限流、按 header key 限流，以及 IP / key 白名单
- 默认配置仍为关闭状态，部署环境需要显式打开

### 启动期告警

- `InitLogging` 现在会在启动时输出基础安全告警
- 当前已覆盖默认或空 `jwt.secret`、开放式 `CORS=*`、关闭限流、空 OpenAI API key、启用 WeKnora 但未配置 API key 等场景
- 当前为 warning-only，不阻断本地开发或测试启动

## 建议的最小生产配置

建议至少开启以下限流基线：

- 全局 `security.rate_limiting.enabled = true`
- `/api/`：保护管理面写操作与批量读取
- `/api/v1/ai/`：单独更严格限流，尤其是 `knowledge/upload`、`knowledge/sync`
- `/api/v1/metrics/ingest`：单独限流，并优先使用 service token 或专用 header key
- `/public/`：按匿名入口单独限流，避免爬取与爆破

建议同时满足以下约束：

- `jwt.secret`、OpenAI / WeKnora API key 不使用代码默认值
- secrets 只通过环境变量或外部 secret manager 注入
- 管理面必须部署在 TLS 后
- `security.cors.allowed_origins` 不使用生产环境全量 `*`
- `X-Request-ID` 由网关或入口统一注入，便于审计追踪

## 仍待补齐的缺口

- token 吊销、refresh token 与 session version 等更细粒度失效机制仍未实现
- 外部开放接口仍缺少代码层强制检查，当前 checklist 主要依赖评审与文档约束
- 关键配置变更尚未形成“执行前审批 / 执行后验证 / 回滚记录”操作手册

## 代码落点

- 路由与表面分类：`apps/server/internal/app/server/router.go`
- 基础中间件装配：`apps/server/internal/app/server/middleware.go`
- 速率限制：`apps/server/internal/middleware/ratelimit.go`
- 审计中间件：`apps/server/internal/platform/audit/gin_middleware.go`
- 启动安全告警：`apps/server/internal/app/bootstrap/security.go`
