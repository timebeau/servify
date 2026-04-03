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

补充说明：

- 当前主业务模型与主运营聚合路径，已基本可以回答“属于哪个 tenant/workspace”
- 仍未完全收口的，主要是少量 legacy 聚合尾项，以及 `DailyStats` 这类系统级全局汇总表是否需要 tenant/workspace 维度拆分的问题

## 核心对象归属现状

### 直接或隐式属于 workspace / tenant 的对象

- `Ticket`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `ticket` 仓储主查询、统计与创建路径已默认按上下文过滤

- `Session`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `conversation` 仓储与 runtime adapter 已默认按上下文过滤

- `Message`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `conversation` 仓储、customer activity、analytics 统计已按上下文过滤

- `TransferRecord`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `routing` 仓储已默认按上下文过滤

- `WaitingRecord`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `routing` 仓储已默认按上下文过滤

- `Customer`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `customer` 仓储通过扩展表 join 收口用户读取、列表、统计与活动查询

- `Agent`
  - 当前已补显式 `tenant_id` / `workspace_id`
  - `agent` 仓储与 transfer runtime load 同步路径已按上下文过滤

- `KnowledgeDoc`
  - 当前 provider 层支持 `tenant_id`
  - 已视为 tenant/workspace scoped 资源，并默认按上下文过滤

- `CustomField`
  - 已补显式 scope 字段
  - service/repository 默认按上下文过滤

- `SLAConfig`
  - 已补显式 scope 字段
  - service 默认按上下文过滤

- `SLAViolation`
  - 已补显式 `tenant_id` / `workspace_id`
  - `sla` service 的违约创建、列表、解决、统计与监控扫描路径已默认按上下文过滤

- `AppIntegration`
  - 已补显式 scope 字段
  - service 默认按上下文过滤

- `ShiftSchedule`
  - 已补显式 `tenant_id` / `workspace_id`
  - `shift` service 的创建、列表、更新、删除与统计路径已默认按上下文过滤

- `Macro`
  - 已补显式 `tenant_id` / `workspace_id`
  - `macro` service 的列表、创建、更新、删除与应用入口已默认按上下文过滤

- `CustomerSatisfaction` / `SatisfactionSurvey`
  - 已补显式 `tenant_id` / `workspace_id`
  - `satisfaction` service 的调查发送、响应、列表、统计、更新与删除路径已默认按上下文过滤

- `SuggestionService`
  - 依赖已 scope 化的 `Ticket` / `KnowledgeDoc`
  - 相似工单与知识库候选查询现已按请求上下文默认过滤，避免跨 workspace 建议结果泄漏

- `GamificationService`
  - 排行榜聚合现已按请求上下文过滤 `Agent` / `Ticket` / `CustomerSatisfaction`
  - 避免跨 workspace 汇总客服绩效与满意度数据

- `MessageRouter` 兼容落库路径
  - 现已从 `UnifiedMessage.Metadata` 提取 `tenant_id` / `workspace_id` 写入 `Session` / `Message`
  - 同一 `sessionID` 若尝试跨 workspace 复用会直接拒绝持久化，避免旧兼容链路串会话

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
- 管理面与服务面受保护路由现已统一拦截显式 scope 选择器：
  - `agent` / `end_user` 不允许在请求中额外声明超出 token 的 `tenant_id` / `workspace_id`
  - `admin` / `service` 仅可在 token 未显式 scope 时通过请求窄化作用域
  - header / query 中若出现互相冲突的 scope 值，直接拒绝请求
- 审计日志会记录这两个字段
- 以下对象已开始显式 scope 化并默认按上下文过滤：
  - `KnowledgeDoc`
  - `CustomField`
  - `SLAConfig`
  - `SLAViolation`
  - `AppIntegration`
  - `ShiftSchedule`
  - `Macro`
  - `CustomerSatisfaction` / `SatisfactionSurvey`
  - `SuggestionService` 的候选查询
- 以下会话链路对象现已补显式 scope 字段，并在 `conversation/routing` 仓储默认按上下文过滤：
  - `Session`
  - `Message`
  - `TransferRecord`
  - `WaitingRecord`
- `Ticket` 现已补显式 scope 字段，并在 `ticket` 仓储的主查询、统计与创建路径按上下文过滤
- `Customer` 现已补显式 scope 字段；`customer` 仓储通过 `customers` 扩展表 join 收口用户读取、列表、统计与活动查询
- `Agent` 现已补显式 scope 字段，并在 `agent` 仓储的创建、读取、列表与统计路径按上下文过滤
- `WorkspaceService`、`analytics` 模块聚合仓储、`gamification` 排行榜聚合、agent transfer load 同步路径已开始按上下文过滤已 scope 化主数据
- 尚未统一 scope 化的主要剩余面：
  - 仍保留在旧 `services/*` 中、未完全迁移到 modules 的少量 legacy 聚合查询
  - `DailyStats` 这类跨租户全局聚合表，当前仍按系统级统计处理，而非 tenant/workspace 维度拆分

## 下一步建议

1. 将 `T4 configuration-scopes` 作为下一阶段重点，明确系统级与租户级配置边界
2. 继续压缩旧 `services/*` 中残留的少量 legacy 聚合查询，避免新增绕过 modules scope 的读写路径
3. 若后续需要租户级报表，补 `DailyStats` 等全局汇总表的 tenant/workspace 维度设计
4. 若要继续降低兼容层风险，逐步把旧 `MessageRouter` / 旧 service 聚合路径迁到 modules 边界
