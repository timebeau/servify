# 11 Tenant Auth And Audit

范围：

- 多租户与 workspace 边界
- 认证授权收口
- 审计日志
- 配置层级治理
- 面向真实部署的安全基线

## T1 tenant-and-workspace-boundaries

- [x] 梳理 tenant、workspace、agent、customer、knowledge base 的归属关系
- [ ] 定义跨租户访问禁止规则
- [x] 定义查询、写入、导出、后台任务的租户隔离语义
- [-] 为核心模型补租户字段与索引策略审查

验收：

- 任何核心数据对象都能明确回答“属于哪个 tenant/workspace”

当前进展：

- 已盘点当前核心对象的 tenant/workspace 归属现状，见 `docs/tenant-workspace-boundaries.md`
- JWT claims 现已支持归一化并透传 `tenant_id` / `workspace_id`
- 认证中间件会将 `tenant_id` / `workspace_id` 投影到请求上下文，供审计、provider 调用和后续数据过滤使用
- `KnowledgeDoc`、`CustomField`、`SLAConfig`、`AppIntegration` 已补显式 scope 字段，并在 service/repository 层默认按上下文过滤
- 当前隔离仍未覆盖全部业务表，`Ticket`、`Session`、`Message`、`Customer`、`Agent` 等主数据仍需后续统一 scope 化

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

- [ ] 区分系统级、租户级、工作区级、运行时级配置
- [ ] 为 AI provider、knowledge provider、routing policy 等配置定义作用域
- [ ] 明确配置加载、覆盖、回退规则
- [ ] 为配置变更补审计与回滚约束

验收：

- 配置不再混杂在环境变量、数据库和代码默认值之间

## T5 security-baseline-for-operations

- [ ] 盘点高风险接口和高风险操作
- [ ] 增加关键操作的速率限制、权限兜底和日志
- [ ] 为 token 生命周期、密钥轮换、敏感字段脱敏补最小规范
- [ ] 为对外开放接口补基础安全清单

验收：

- 项目具备进入真实部署前的最小安全治理骨架
