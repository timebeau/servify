# 10 Service To Module Migration

范围：

- 旧 `services` / `handlers` 收口
- `modules/*` 应用入口固化
- 兼容适配层
- 迁移路线与完成度治理

## M1 migration-inventory

- [x] 盘点 `handlers -> services -> repo` 旧链路
- [x] 盘点 `delivery -> application -> domain -> infra` 新链路覆盖情况
- [x] 为每个领域模块绘制当前入口与目标入口映射
- [x] 标出高风险交叉依赖、共享模型、循环引用点

验收：

- 每个核心能力都能回答“当前走哪条链路，目标走哪条链路”

产出：

- [10-migration-inventory.md](./10-migration-inventory.md)

## M2 module-entry-contracts

- [-] 为 `conversation`、`routing`、`ticket`、`analytics`、`ai` 定义唯一应用入口
- [-] 明确 handler 只能依赖 adapter / application contract，不能直接拼业务
- [-] 为旧 service 定义允许保留的兼容职责
- [ ] 为 repo / gorm 依赖建立单向约束

验收：

- 新增需求不会再同时落入旧 service 和新 module 两套实现

当前已落地：

- `ticket`
  - handler-facing contract 已收口到 `modules/ticket/delivery.HandlerService`
  - runtime 已使用 `HandlerServiceAdapter`
  - `services/TicketService` 明确保留为兼容 facade，而不是 HTTP 入口
- `agent`
  - handler-facing contract 已收口到 `modules/agent/delivery.HandlerService`
  - runtime 仍复用 `services.AgentService` 作为 contract 实现
  - `services.AgentService` 明确保留为兼容 facade + runtime state holder
- `analytics`
  - handler-facing contract 已收口到 `modules/analytics/delivery.HandlerService`
  - runtime 仍复用 `services.StatisticsService` 作为 contract 实现
  - `services.StatisticsService` 明确保留为兼容 facade + event bus glue
- `routing / session transfer`
  - handler-facing contract 已收口到 `modules/routing/delivery.HandlerService`
  - runtime 仍复用 `services.SessionTransferService` 作为 contract 实现
  - `services.SessionTransferService` 明确保留为兼容 facade + runtime glue
- `conversation / websocket runtime`
  - websocket 持久化入口已显式收口到 `modules/conversation/delivery.WebSocketMessageWriter`
  - 主 runtime 与 lightweight realtime runtime 都通过 `WebSocketMessageAdapter` 注入 conversation 模块
- `ai`
  - handler-facing contract 已收口到 `modules/ai/delivery.HandlerService`
  - `AIAssembly` 已拆分为 handler-facing service 与 runtime-facing service
  - 默认 handler 主路径已提升到 orchestrated enhanced AI，而不是直接暴露旧 `AIServiceInterface`
- 边界守护
  - CI 已增加 `scripts/check-module-boundaries.sh`
  - 当前自动校验 `ticket` / `agent` / `analytics` / `routing` / `ai` 的 handler constructor、router dependency、runtime wiring，以及 `conversation` 的 websocket persistence wiring，都必须停留在 module delivery contract

## M3 compatibility-adapters

- [ ] 为仍依赖旧 service 的 handler 增加过渡 adapter
- [ ] 将旧 service 中可迁移的业务规则搬到 module application 层
- [ ] 将旧 service 缩为兼容 facade，而不是业务中心
- [ ] 为适配层增加测试，确保迁移前后行为一致

验收：

- 可以逐步替换旧链路，而不需要一次性重写

## M4 core-module-migrations

- [ ] 优先迁移 `conversation`
- [ ] 优先迁移 `routing`
- [ ] 优先迁移 `ticket`
- [ ] 优先迁移 `ai`
- [ ] 为每个模块拆分 1 到 3 个可独立提交的小任务包

验收：

- 核心闭环模块率先完成收口，后续模块可复用迁移模式

## M5 boundary-enforcement-and-scorecard

- [ ] 增加迁移完成度表
- [ ] 为目录职责与依赖方向补文档
- [ ] 在 code review / CI 中增加边界约束检查
- [ ] 明确 `services/*` 的冻结策略与退役条件

验收：

- 迁移工作可追踪、可量化，不会长期停留在“正在重构”
