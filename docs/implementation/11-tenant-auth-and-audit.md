# 11 Tenant Auth And Audit

范围：

- 多租户与 workspace 边界
- 认证授权收口
- 审计日志
- 配置层级治理
- 面向真实部署的安全基线

## T1 tenant-and-workspace-boundaries

- [ ] 梳理 tenant、workspace、agent、customer、knowledge base 的归属关系
- [ ] 定义跨租户访问禁止规则
- [ ] 定义查询、写入、导出、后台任务的租户隔离语义
- [ ] 为核心模型补租户字段与索引策略审查

验收：

- 任何核心数据对象都能明确回答“属于哪个 tenant/workspace”

## T2 auth-and-rbac-convergence

- [ ] 盘点当前 JWT、claims、permissions、middleware 的实际入口
- [ ] 收口 RBAC 模型与权限解析链路
- [ ] 区分 end-user、agent、admin、service token 的认证语义
- [ ] 为管理端、公开接口、内部接口定义不同授权策略
- [ ] 为关键权限补负向测试与越权测试

验收：

- 权限判断路径单一且可测试，不依赖散落的 handler 逻辑

## T3 audit-log-foundation

- [ ] 定义审计事件模型
- [ ] 覆盖关键写操作，例如工单变更、路由分配、配置变更、权限变更
- [ ] 记录 actor、tenant、resource、before/after、request metadata
- [ ] 设计查询接口与保留策略
- [ ] 为敏感操作提供最小可追溯能力

验收：

- 关键变更可追溯，满足问题排查和合规基础需求

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

