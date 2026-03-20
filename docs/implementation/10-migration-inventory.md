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
  - 但主流程仍强依赖 `gorm.DB`、`AgentService`、`WebSocketHub` 和旧会话模型
  - 结论：属于“局部接入 module adapter，但主流程仍是 legacy service”的高风险混合区

- `ai`
  - 同时存在旧 `services/AIService`、`EnhancedAIService`、`OrchestratedAIService`
  - `OrchestratedAIService` 已是 `modules/ai/application.QueryOrchestrator` 的适配层
  - handler 仍面向 `AIServiceInterface`，但接口下挂着多套实现
  - 结论：迁移重点不是单纯改 handler，而是先收敛“哪一个才是默认主路径”

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

这些模块的共同特点：

- module 已存在
- 但核心运行时仍由旧 service 或旧 runtime struct 主导
- 需要先画清 runtime 职责边界，再做 handler/adapter 收口

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
4. `routing / session transfer`
5. `ai`

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
- `routing / session transfer`
  - HTTP handler 入口：`modules/routing/delivery.HandlerService`
  - 旧 `services/SessionTransferService` 定位：兼容 facade + runtime glue + legacy transfer orchestration

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
- `services/SessionTransferService`
  - 允许保留：旧调用方兼容入口、AgentService/WebSocketHub/AIService 协调、转接实时通知与 legacy transfer orchestration
  - 不应新增：新的 HTTP handler 直接依赖、等待队列之外继续扩散 routing 读写规则

## 当前自动化守护

- `scripts/check-module-boundaries.sh`
  - 校验 `ticket` / `agent` / `analytics` / `routing` 的 handler constructor 必须依赖 `modules/*/delivery.HandlerService`
  - 校验 router/runtime 对这四个模块的注入类型必须停留在 handler-facing contract
  - 目的：先锁住已完成迁移的入口，避免回退到 handler 直连具体旧 service
