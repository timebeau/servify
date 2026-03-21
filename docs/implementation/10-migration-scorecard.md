# 10 Migration Scorecard

本文件记录 `10-service-to-module-migration` 的当前完成度，用于支持 M5 的持续追踪。

状态定义：

- `legacy`
  - handler/runtime 仍主要依赖旧 `services/*`
- `mixed`
  - 已部分接入 `modules/*`，但 HTTP 或 runtime 主入口尚未收口
- `contracted`
  - handler-facing 或 runtime-facing 主入口已收口到 `modules/*/delivery`
- `stabilized`
  - 已收口，且有 CI / 文档 / 兼容职责规则持续守护

判定补充：

- 若某能力已经收口，但尚未进入 `scripts/module-boundaries.rules`，不能标记为 `stabilized`
- 若文档已声明收口，但 CI 规则尚未覆盖，最多只能标记为 `contracted`

## Scorecard

| Capability | Handler Entry | Runtime Entry | Legacy Service Role | Status | Notes |
| --- | --- | --- | --- | --- | --- |
| `ticket` | `modules/ticket/delivery.HandlerService` | `ticketdelivery.NewHandlerServiceAdapter(...)` | compatibility facade + orchestration side effects | `stabilized` | handler/router/runtime 已收口，边界脚本已覆盖 |
| `agent` | `modules/agent/delivery.HandlerService` | `services.AgentService` as contract impl | compatibility facade + runtime state holder | `stabilized` | 运行态状态仍在 legacy service，但 HTTP 入口已锁定 |
| `analytics` | `modules/analytics/delivery.HandlerService` | `services.StatisticsService` as contract impl | compatibility facade + event bus glue | `stabilized` | DTO 已抽到 module contract，handler 不再直连 concrete service |
| `routing / session transfer` | `modules/routing/delivery.HandlerService` | `services.SessionTransferService` as contract impl | compatibility facade + runtime glue | `stabilized` | waiting queue 与 transfer record 写入已统一走 module，主流程仍保留 session/ticket/ws 协调 |
| `conversation / websocket runtime` | n/a | `modules/conversation/delivery.WebSocketMessageWriter` | runtime connection hub | `stabilized` | 主 runtime 与 lightweight realtime runtime 都已走 adapter 注入 |
| `ai` | `modules/ai/delivery.HandlerService` | `AIAssembly.RuntimeService` | runtime compatibility surface for websocket/router/transfer | `stabilized` | handler 主路径已切到 orchestrated enhanced AI |
| `customer` | `modules/customer/delivery.HandlerService` | `services.CustomerService` as contract impl | compatibility facade + DTO mapping | `contracted` | handler/router 已收口，runtime 仍保留 legacy facade |
| `automation` | `modules/automation/delivery.HandlerService` | `services.AutomationService` as contract impl | compatibility facade + event bus glue | `contracted` | handler/router 已收口，runtime 仍保留 subscriber / event bus glue |
| `knowledge` | `modules/knowledge/delivery.HandlerService` | `services.KnowledgeDocService` as contract impl | compatibility facade + module assembly | `contracted` | handler/router 已收口，index job runtime 仍未接入 |
| `voice` | already module coordinator | module coordinator | not primarily legacy | `mixed` | 已较模块化，但不属于本轮 `services -> modules` 的典型迁移样式 |

## 当前结论

- `ticket`、`agent`、`analytics`、`routing`、`conversation runtime`、`ai` 已具备持续守护条件
- `customer` 已完成 handler-facing contract 收口，但 runtime 仍保留 legacy facade
- `automation` 已完成 handler-facing contract 收口，但 runtime 仍保留 event bus subscriber 与旧兼容 facade
- `knowledge` 已完成 handler-facing contract 收口，但索引任务与 provider 集成仍未纳入 runtime 主路径
- 上述 `stabilized` 条目都已进入 `scripts/module-boundaries.rules`
- 旧 `services/*` 仍然存在，但对这些能力来说，已不再是 HTTP 默认入口
- 下一阶段重点不再是“再找一个模块收口”，而是：
  - 扩大 scorecard 覆盖范围
  - 持续压缩剩余 `legacy` 状态条目，并补齐 `knowledge` 的 indexing/runtime 策略
  - 为 `legacy` 条目明确冻结策略与退役条件
