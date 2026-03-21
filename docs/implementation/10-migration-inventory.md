# 10 Migration Inventory

本文件记录 `handlers -> services -> modules` 的当前迁移现状，用于支持 `10-service-to-module-migration` 的 M1 盘点阶段。

## 当前观察结论

- `ticket`
  - 已具备较强的模块化基础
  - `handlers/ticket_handler.go` 通过 `ticketHandlerService` 接口工作
  - 已存在模块交付适配层 `modules/ticket/delivery/handler_adapter.go`
  - 但旧 `services/TicketService` 仍保留，并且本身也在包装 `ticket` module 与 orchestration
  - 当前规则：HTTP handler 应依赖 `modules/ticket/delivery.HandlerService`
  - 结论：这是最适合优先完成“handler 直接面向 module adapter”的模块

- `agent`
  - `services/AgentService` 已明显成为兼容层
  - 核心读写已走 `modules/agent/application.Service`
  - 但仍保留 `onlineAgents`、`agentQueues`、后台同步等 legacy runtime 状态
  - 当前规则：HTTP handler 应依赖 `modules/agent/delivery.HandlerService`
  - 结论：属于“module 已落地，但 runtime 兼容状态尚未收口”的类型

- `analytics`
  - `services/StatisticsService` 基本是 module facade
  - 核心统计能力已走 `modules/analytics/application.Service`
  - 还保留事件总线订阅和 DTO 映射的兼容包装
  - 结论：迁移风险低，适合作为后续收口样板

- `routing / session transfer`
  - `services/SessionTransferService` 已注入 `modules/routing/delivery.SessionTransferAdapter`
  - waiting queue 的读取、新增、查询、取消、转态同步，以及 transfer record 写入已收口到 routing module 状态机
  - 但主流程仍强依赖 `gorm.DB`、`AgentService`、`WebSocketHub` 和旧会话模型
  - 结论：属于“局部接入 module adapter，但主流程仍是 legacy service”的高风险混合区

- `ai`
  - 同时存在旧 `services/AIService`、`EnhancedAIService`、`OrchestratedAIService`
  - `OrchestratedAIService` 已是 `modules/ai/application.QueryOrchestrator` 的适配层
  - handler 已开始收口到 `modules/ai/delivery.HandlerService`
  - 结论：迁移重点是继续压缩 legacy `AIServiceInterface` 的 handler 可见面，只保留 runtime 兼容用途

- `customer`
  - `services/CustomerService` 已经是 `modules/customer/application.Service` 的轻量 facade
  - HTTP handler 只消费请求/响应 DTO 与兼容方法，核心业务已经下沉到 module application + infra repository
  - 结论：适合先定义 handler-facing contract，并把 router/runtime 注入切到 `modules/customer/delivery`

- `automation`
  - `services/AutomationService` 已经把核心触发器与执行查询下沉到 `modules/automation/application.Service`
  - legacy service 主要保留 event bus subscriber 注册、测试兼容 helper 与 module 装配
  - 结论：适合先定义 handler-facing contract，并保留 runtime/event bus glue 在旧 service

- `knowledge`
  - 已存在 `modules/knowledge/application`，但此前缺少 Gorm repository 与 handler-facing delivery contract
  - 旧 `services/KnowledgeDocService` 主要承接 CRUD 与 tag/category DTO，适合改成 module facade
  - 结论：先补 module 落地层和 handler contract，再把 runtime 中的 provider/indexing 接口留待后续收口

## 按迁移成熟度分组

### A. 已有明显 module facade

- `ticket`
- `agent`
- `analytics`

这些模块的共同特点：

- 已有 `modules/*/application` 作为核心业务入口
- 旧 `services/*` 更多承担兼容包装、DTO 映射、runtime glue
- 适合优先定义“唯一入口”并冻结旧 service 新增逻辑

### B. 已局部接入 module adapter，但主流程仍偏旧

- `routing / session transfer`
- `conversation` 相关 websocket/runtime 路径
- `customer`

这些模块的共同特点：

- module 已存在
- 但核心运行时仍由旧 service 或旧 runtime struct 主导
- 需要先画清 runtime 职责边界，再做 handler/adapter 收口
- `customer` 的 handler 已可先行收口，但 runtime 仍保留旧 facade 用于兼容 DTO 和构造装配
- `automation` 的 handler 已可先行收口，但 runtime 仍需要旧 service 持有 event bus subscriber
- `knowledge` 的 handler 已可先行收口，但 indexing/provider 尚未进入统一 runtime 装配

### C. 多实现并存，需要先确定默认主路径

- `ai`

特点：

- legacy AI service 与 orchestrated AI module 并存
- provider、fallback、enhanced 包装已经很多
- 下一步更像“架构收敛”而不是简单替换调用点

## 推荐迁移顺序

1. `ticket`
2. `agent`
3. `analytics`
4. `customer`
5. `automation`
6. `knowledge`
7. `routing / session transfer`
8. `ai`

## 当前高风险点

- 同一个业务能力同时存在 handler interface、legacy service、module service、delivery adapter 多层入口
- 旧 service 中仍混有 runtime 状态与 side effects，导致无法简单替换
- `ai` 路径下存在多个“看起来都像主实现”的对象，容易继续分叉
- `routing` 与 `conversation` 仍强依赖旧 websocket / session runtime

## 下一步建议

- 先在 `ticket` 上定义唯一 HTTP 入口应该直连哪个 adapter
- 明确 `services/*` 中哪些类型属于：
  - 兼容 facade
  - runtime state holder
  - 仍未迁出的真实业务层
- 为 `ticket`、`agent`、`analytics` 各拆一个最小任务包

## 已确认的入口规则

- `ticket`
  - HTTP handler 入口：`modules/ticket/delivery.HandlerService`
  - 旧 `services/TicketService` 定位：兼容 facade + orchestration side effects 组装
- `agent`
  - HTTP handler 入口：`modules/agent/delivery.HandlerService`
  - 旧 `services/AgentService` 定位：兼容 facade + runtime state holder
- `analytics`
  - HTTP handler 入口：`modules/analytics/delivery.HandlerService`
  - 旧 `services/StatisticsService` 定位：兼容 facade + event bus subscriber / DTO mapping glue
- `customer`
  - HTTP handler 入口：`modules/customer/delivery.HandlerService`
  - 旧 `services/CustomerService` 定位：兼容 facade + DTO mapping + module application 装配
- `automation`
  - HTTP handler 入口：`modules/automation/delivery.HandlerService`
  - 旧 `services/AutomationService` 定位：兼容 facade + event bus subscriber + module application 装配
- `knowledge`
  - HTTP handler 入口：`modules/knowledge/delivery.HandlerService`
  - 旧 `services/KnowledgeDocService` 定位：兼容 facade + module application / repository 装配
- `routing / session transfer`
  - HTTP handler 入口：`modules/routing/delivery.HandlerService`
  - 旧 `services/SessionTransferService` 定位：兼容 facade + runtime glue + legacy transfer orchestration
- `conversation / websocket runtime`
  - websocket 持久化入口：`modules/conversation/delivery.WebSocketMessageWriter`
  - `services/WebSocketHub` 定位：runtime connection hub，不再自定义 conversation 私有持久化接口
- `ai`
  - HTTP handler 入口：`modules/ai/delivery.HandlerService`
  - `AIAssembly` 定位：显式区分 handler-facing AI contract 与 runtime-facing `AIServiceInterface`

## 已确认的兼容职责边界

- `services/TicketService`
  - 允许保留：旧调用方兼容入口、event bus / automation / satisfaction 等 side effect 组装、module command/orchestrator 的装配
  - 不应新增：新的 HTTP handler 直接依赖、新业务规则主入口、绕过 `modules/ticket/delivery.HandlerService` 的写路径
- `services/AgentService`
  - 允许保留：旧调用方兼容入口、在线客服运行时状态、队列与后台同步逻辑
  - 不应新增：新的 HTTP handler 直接依赖、脱离 `modules/agent/application.Service` 的核心业务写路径
- `services/StatisticsService`
  - 允许保留：旧调用方兼容入口、event bus 订阅注册、DTO 映射、统计后台任务调度
  - 不应新增：新的 HTTP handler 直接依赖、绕过 `modules/analytics/application.Service` 的主统计读写路径
- `services/CustomerService`
  - 允许保留：旧调用方兼容入口、DTO 映射、`modules/customer/application.Service` 与 repository 的装配
  - 不应新增：新的 HTTP handler 直接依赖、绕过 `modules/customer/application.Service` 的核心客户读写路径
- `services/AutomationService`
  - 允许保留：旧调用方兼容入口、event bus subscriber 注册、测试辅助方法、`modules/automation/application.Service` 的装配
  - 不应新增：新的 HTTP handler 直接依赖、绕过 `modules/automation/application.Service` 的触发器和执行记录主入口
- `services/KnowledgeDocService`
  - 允许保留：旧调用方兼容入口、DTO 映射、`modules/knowledge/application.Service` 与 repository 的装配
  - 不应新增：新的 HTTP handler 直接依赖、绕过 `modules/knowledge/application.Service` 的文档 CRUD 主入口
- `services/SessionTransferService`
  - 允许保留：旧调用方兼容入口、AgentService/WebSocketHub/AIService 协调、转接实时通知与 legacy transfer orchestration
  - 不应新增：新的 HTTP handler 直接依赖、绕过 routing module 的 waiting/transfer record 读写规则
- `services/WebSocketHub`
  - 允许保留：连接管理、广播、协议消息分发、对 AI/transfer 的 runtime glue
  - 不应新增：新的 conversation 私有持久化接口、绕过 `modules/conversation/delivery.WebSocketMessageWriter` 的消息落库路径
- `services/AIServiceInterface` 及 legacy AI types
  - 允许保留：WebSocketHub、MessageRouter、SessionTransferService 等 runtime 兼容调用面
  - 不应新增：新的 HTTP handler 直接依赖、继续扩散为默认 AI 主入口

## 当前自动化守护

- `scripts/check-module-boundaries.sh`
  - 校验 `ticket` / `agent` / `analytics` / `customer` / `automation` / `knowledge` / `routing` / `ai` 的 handler constructor 必须依赖 `modules/*/delivery.HandlerService`
  - 校验 router/runtime 对这八个模块的注入类型必须停留在 handler-facing contract，并校验 `conversation` 的 websocket persistence 入口必须走 module delivery adapter
  - 目的：先锁住已完成迁移的入口，避免回退到 handler 直连具体旧 service
