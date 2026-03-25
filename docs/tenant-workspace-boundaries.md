# Tenant And Workspace Boundaries

本文件记录当前仓库中核心对象的归属关系与隔离语义，作为 `11 / T1 tenant-and-workspace-boundaries` 的基础盘点。

## 当前最小归属模型

- `tenant`
  - 代表部署级业务租户边界
  - 当前已在配置与 provider 参数中出现，例如 WeKnora `tenant_id`
  - 当前业务表大多尚未显式落库 `tenant_id`

- `workspace`
  - 代表租户下的操作空间或管理空间
  - 当前已在权限语义中出现 `workspace.read`
  - 当前仓库存在工作台概览能力，但尚未形成独立 `workspace` 持久化模型

## 核心对象归属现状

### 直接或隐式属于 workspace / tenant 的对象

- `Ticket`
  - 当前通过 `Customer`、`Agent`、`Session` 等关系隐式归属
  - 尚无显式 `tenant_id` / `workspace_id`

- `Session`
  - 当前与 `User` / `Agent` / `Ticket` 关联
  - 归属仍为隐式

- `Message`
  - 当前通过 `Session` 继承归属

- `TransferRecord`
  - 当前通过 `SessionID` 继承归属

- `WaitingRecord`
  - 当前通过 `SessionID` 继承归属

- `KnowledgeDoc`
  - 当前 provider 层支持 `tenant_id`
  - 应视为 tenant-scoped 资源

- `CustomField`
  - 当前更接近 workspace-scoped 配置
  - 但尚未显式建模作用域

- `SLAConfig`
  - 当前更接近 workspace-scoped 或 tenant-scoped 配置
  - 尚未显式建模

- `AppIntegration`
  - 当前更接近 workspace-scoped 配置
  - 尚未显式建模

## 当前运行时上下文

JWT claims 已支持透传：

- `tenant_id`
- `workspace_id`

并在认证中间件中投影到请求上下文，供后续：

- 数据访问过滤
- 审计日志补维
- provider 调用透传
- 后台任务隔离

同时，这两个字段现在也会进入 `request context`，因此 service / repository 层可以直接读取，而不必依赖 gin transport context。

## 当前隔离语义

- 管理面与服务面请求已经可以带 `tenant_id` / `workspace_id` 上下文
- 审计日志会记录这两个字段
- 以下对象已开始显式 scope 化并默认按上下文过滤：
  - `KnowledgeDoc`
  - `CustomField`
  - `SLAConfig`
  - `AppIntegration`
- 其余业务表读写尚未统一按 `tenant_id` / `workspace_id` 过滤

## 下一步建议

1. 为真正 tenant-scoped 的核心表补显式 `tenant_id`
2. 为 workspace-scoped 配置表补 `workspace_id`
3. 先从配置类与知识库类对象开始，因为它们的边界最清晰
4. 在 repository / service 层引入统一 scope 过滤，而不是把过滤散落在 handler
