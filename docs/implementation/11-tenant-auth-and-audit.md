# 11 Tenant Auth And Audit

范围：

- 多租户与 workspace 边界
- 认证授权收口
- 审计日志
- 配置层级治理
- 面向真实部署的安全基线

## T1 tenant-and-workspace-boundaries

- [x] 梳理 tenant、workspace、agent、customer、knowledge base 的归属关系
- [x] 定义跨租户访问禁止规则
- [x] 定义查询、写入、导出、后台任务的租户隔离语义
- [x] 为核心模型补租户字段与索引策略审查

验收：

- 任何核心数据对象都能明确回答“属于哪个 tenant/workspace”

当前进展：

- 已盘点当前核心对象的 tenant/workspace 归属现状，见 `docs/tenant-workspace-boundaries.md`
- JWT claims 现已支持归一化并透传 `tenant_id` / `workspace_id`
- 认证中间件会将 `tenant_id` / `workspace_id` 投影到请求上下文，供审计、provider 调用和后续数据过滤使用
- 管理面 `/api/*` 与受保护的 `/api/v1/*` 已增加统一 request-scope 守卫，拒绝 token scope 与 header/query scope 冲突的请求，并限制非 admin/service principal 通过请求扩大作用域
- `KnowledgeDoc`、`CustomField`、`SLAConfig`、`AppIntegration` 已补显式 scope 字段，并在 service/repository 层默认按上下文过滤
- `Session`、`Message`、`TransferRecord`、`WaitingRecord` 已补显式 `tenant_id/workspace_id` 字段与基础索引，并在 `conversation/routing` 仓储默认按上下文过滤
- `Ticket` 已补显式 `tenant_id/workspace_id` 字段与基础索引，并在 `ticket` 仓储主查询、统计与创建路径按上下文过滤
- `Customer` 已补显式 `tenant_id/workspace_id` 字段，并在 `customer` 仓储读取、列表、统计与活动查询通过扩展表 scope 收口
- `Agent` 已补显式 `tenant_id/workspace_id` 字段，并在 `agent` 仓储创建、读取、列表与统计路径按上下文过滤
- `WorkspaceService`、`analytics` 模块聚合仓储与 agent transfer runtime load 同步路径已开始按上下文过滤 `Session/Agent/Ticket/Message/Customer` 等已 scope 化主数据
- 当前剩余缺口已从“核心主数据未建模”收敛到“部分运营模型与旧 service 聚合查询未继续 scope 化”，详见 `docs/tenant-workspace-boundaries.md`

## T2 auth-and-rbac-convergence

- [x] 盘点当前 JWT、claims、permissions、middleware 的实际入口
- [x] 收口 RBAC 模型与权限解析链路
- [x] 区分 end-user、agent、admin、service token 的认证语义
- [x] 为管理端、公开接口、内部接口定义不同授权策略
- [x] 为关键权限补负向测试与越权测试

验收：

- 权限判断路径单一且可测试，不依赖散落的 handler 逻辑

当前进展：

- JWT 校验、claims 归一化、role -> permission 展开、principal kind 判定已经统一收口到 `internal/platform/auth`
- `internal/middleware` 仅保留兼容入口，不再承载独立权限逻辑
- 管理面 `/api/*` 已显式拒绝 `end_user` token，当前仅允许 `agent`、`admin`、`service`
- `public` 路由保持匿名访问；`/api/v1` 已拆分为 public realtime、management realtime/AI、service-only metrics ingest 三类表面
- 表面分类与新增路由准入规则已文档化，见 `docs/auth-surface-policy.md`

## T3 audit-log-foundation

- [x] 定义审计事件模型
- [x] 覆盖关键写操作，例如工单变更、路由分配、配置变更、权限变更
- [x] 记录 actor、tenant、resource、before/after、request metadata
- [x] 设计查询接口与保留策略
- [x] 为敏感操作提供最小可追溯能力

验收：

- 关键变更可追溯，满足问题排查和合规基础需求

当前进展：

- 已新增 `AuditLog` 持久化模型，并纳入应用迁移集合
- 管理面 `/api/*` 写请求已统一接入审计中间件，成功写操作会记录 actor、principal、action、resource、request metadata 与请求载荷
- `before/after` 已预留中间件扩展点，后续可为高价值 handler 补充精细状态快照
- 已新增 `GET /api/audit/logs` 查询接口，支持按 action/resource/principal/actor/success/time-range 过滤
- 审计查询与保留基线已文档化，见 `docs/audit-log-policy.md`
- 当前仍未接入归档/清理 worker，保留策略暂为文档约束而非自动执行

## T4 configuration-scopes

- [x] 区分系统级、租户级、工作区级、运行时级配置
- [x] 为 AI provider、knowledge provider、routing policy 等配置定义作用域
- [x] 明确配置加载、覆盖、回退规则
- [x] 为配置变更补审计与回滚约束

验收：

- 配置不再混杂在环境变量、数据库和代码默认值之间

当前进展：

- 已新增 `docs/configuration-scopes.md`，明确 system / tenant / workspace / runtime 四层配置作用域
- 已定义 AI provider、knowledge provider、routing policy、portal、security baseline 的推荐作用域矩阵
- 已定义统一覆盖顺序：`runtime -> workspace -> tenant -> system -> code default`
- 已定义配置变更的审计与回滚约束，区分系统级发布回滚与租户/工作区配置恢复
- 已新增 `internal/platform/configscope` 作为统一 resolver 骨架，当前已覆盖 portal config、OpenAI config 与 WeKnora config 的分层解析
- 门户公开配置读取、AI runtime 装配与增强健康检查中的 WeKnora 展示信息已切到 resolver，不再各自拼接 system default / runtime override
- 当前仍未实现通用 `TenantConfig` / `WorkspaceConfig` 持久化 provider；租户 / 工作区级覆盖规则已文档化，但代码层仍以 system config + 少量 scoped data path 为主

## T5 security-baseline-for-operations

- [ ] 盘点高风险接口和高风险操作
- [ ] 增加关键操作的速率限制、权限兜底和日志
- [ ] 为 token 生命周期、密钥轮换、敏感字段脱敏补最小规范
- [ ] 为对外开放接口补基础安全清单

验收：

- 项目具备进入真实部署前的最小安全治理骨架

当前进展：

- 已新增 `docs/security-baseline-operations.md`，盘点 management / service / public 三类高风险表面与关键接口
- 已确认基础路由统一挂载 `RateLimitMiddlewareFromConfig`，具备全局、路径前缀与 header key 级限流能力
- 已确认管理面写接口统一叠加 `AuthMiddleware`、`EnforceRequestScope`、`RequirePrincipalKinds`、`RequireResourcePermission` 与 `AuditMiddleware`
- 审计中间件现已对 `password`、`secret`、`api_key`、`token` 等敏感字段做统一脱敏，避免明文 secrets 落入审计日志
- 已新增 `docs/token-lifecycle-and-key-rotation.md`，定义 JWT / provider key / service token 的最小生命周期与轮换流程
- 已新增 `config.production.secure.example.yml`，提供启用限流、收紧 CORS、使用环境变量注入 secrets 的生产示例模板
- 已新增 `docs/public-surface-security-checklist.md`，收口匿名 / 对外开放接口的最小评审项与逐类入口检查清单
- 已新增 bootstrap 启动安全告警，对默认 `jwt.secret`、开放式 CORS、关闭限流、空 provider key 等高风险配置输出 warning
- 启动安全告警现已进一步覆盖匿名入口限流基线，若 `/public/` 或 `/api/v1/ws` 未配置独立路径级 rate limit，会在启动时明确告警
- `internal/platform/auth` 现已新增可组合 token policy 钩子，支持按 `iat` 做 issued-before 失效、按 `token_version` 做最小会话版本淘汰，为后续账号级吊销策略预留接入点
- `internal/platform/auth` 现已新增基于 `users` 表状态的 token policy，并接入 router auth middleware：非 active 用户会被拒绝，设置 `token_valid_after` / `token_version` 后可使旧 token 失效
- agent 管理面现已新增 `POST /api/agents/:id/revoke-tokens`，可主动提升用户 `token_version` 并刷新 `token_valid_after`，将旧 token 作废
- customer 管理面现已新增 `POST /api/customers/:id/revoke-tokens`，可对受作用域约束的客户账号主动触发旧 token 失效
- token 主动失效底层逻辑现已收敛到 `internal/platform/usersecurity`，agent/customer 两条管理面路径复用同一套 `token_valid_after + token_version` 更新实现
- 管理面现已新增统一 `security` surface：`GET /api/security/users/:id` 可查询用户安全状态，`POST /api/security/users/:id/revoke-tokens` 可执行通用 user revoke，分别由 `security.read` / `security.write` 权限保护
- 当前仍缺 refresh token / revoke list 等更细粒度失效实现，以及批量失效、审批回滚等更完整的 user security 操作能力
