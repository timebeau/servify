# 03 Business Modules

范围：

- conversation
- routing
- agent
- ticket
- customer
- automation
- analytics
- voice

## B1 conversation

- [ ] 创建 `internal/modules/conversation/domain`
- [ ] 创建 `internal/modules/conversation/application`
- [ ] 创建 `internal/modules/conversation/infra`
- [ ] 创建 `internal/modules/conversation/delivery`
- [ ] 定义 `Conversation`
- [ ] 定义 `Participant`
- [ ] 定义 `ConversationMessage`
- [ ] 定义 `ChannelBinding`
- [ ] 定义 `ConversationStatus`
- [ ] 实现 CreateConversation
- [ ] 实现 ResumeConversation
- [ ] 实现 IngestTextMessage
- [ ] 实现 IngestSystemEvent
- [ ] 发布 `conversation.created`
- [ ] 发布 `conversation.message_received`

验收：

- conversation 成为统一交互模型

## B2 routing

- [ ] 创建 `internal/modules/routing/domain`
- [ ] 创建 `internal/modules/routing/application`
- [ ] 创建 `internal/modules/routing/infra`
- [ ] 创建 `internal/modules/routing/delivery`
- [ ] 定义 `Assignment`
- [ ] 定义 `QueueEntry`
- [ ] 定义 `TransferRecord`
- [ ] 定义 `AgentAvailabilityPolicy`
- [ ] 定义 `SkillMatchPolicy`
- [ ] 定义 `LoadBalancePolicy`
- [ ] 实现 RequestHumanHandoff
- [ ] 实现 AssignAgent
- [ ] 实现 AddToWaitingQueue
- [ ] 实现 CancelWaiting
- [ ] 发布 `routing.agent_assigned`
- [ ] 发布 `routing.transfer_completed`

验收：

- 取代现有 session transfer 交叉依赖

## B3 agent

- [ ] 创建 `internal/modules/agent/domain`
- [ ] 创建 `internal/modules/agent/application`
- [ ] 创建 `internal/modules/agent/infra`
- [ ] 创建 `internal/modules/agent/delivery`
- [ ] 定义 `AgentProfile`
- [ ] 定义 `AgentPresence`
- [ ] 定义 `AgentLoad`
- [ ] 增加 chat concurrency
- [ ] 增加 voice concurrency
- [ ] 实现 GoOnline
- [ ] 实现 GoOffline
- [ ] 实现 MarkBusy
- [ ] 实现 MarkAway

验收：

- agent 模块支持未来 voice load

## B4 ticket

- [x] 创建 `internal/modules/ticket/domain`
- [x] 创建 `internal/modules/ticket/application`
- [x] 创建 `internal/modules/ticket/infra`
- [x] 创建 `internal/modules/ticket/delivery`
- [x] 新增 `query_service.go`
- [x] 迁移 `ListTickets`
- [x] 迁移 `GetTicketByID`
- [x] 新增 `command_service.go`
- [x] 迁移 `CreateTicket`
- [x] 迁移 `UpdateTicket`
- [x] 迁移 `AssignTicket`
- [x] 迁移 `AddComment`
- [x] 迁移 `CloseTicket`
- [x] 拆 `StatusTransitionPolicy`
- [x] 拆 `CustomFieldValidator`
- [ ] 拆 `BulkUpdateTickets`
- [ ] 发布 `ticket.created`
- [ ] 发布 `ticket.assigned`
- [ ] 发布 `ticket.closed`

验收：

- 旧 TicketService 不再是超大单体

## B5 customer

- [ ] 创建 `internal/modules/customer/domain`
- [ ] 创建 `internal/modules/customer/application`
- [ ] 创建 `internal/modules/customer/infra`
- [ ] 创建 `internal/modules/customer/delivery`
- [ ] 定义 `CustomerProfile`
- [ ] 定义 `CustomerTag`
- [ ] 定义 `CustomerNote`
- [ ] 定义 `CustomerActivity`
- [ ] 实现 CreateCustomer
- [ ] 实现 UpdateCustomer
- [ ] 实现 AddNote
- [ ] 实现 UpdateTags
- [ ] 实现 GetCustomerActivity

验收：

- customer 模块独立成型

## B6 automation

- [ ] 创建 `internal/modules/automation/domain`
- [ ] 创建 `internal/modules/automation/application`
- [ ] 创建 `internal/modules/automation/infra`
- [ ] 创建 `internal/modules/automation/delivery`
- [ ] 定义 `Trigger`
- [ ] 定义 `Condition`
- [ ] 定义 `Action`
- [ ] 定义 `Execution`
- [ ] 订阅 ticket events
- [ ] 订阅 conversation events
- [ ] 订阅 routing events
- [ ] 迁移 dry-run

验收：

- automation 基于 event bus 执行

## B7 analytics

- [ ] 创建 `internal/modules/analytics/domain`
- [ ] 创建 `internal/modules/analytics/application`
- [ ] 创建 `internal/modules/analytics/infra`
- [ ] 创建 `internal/modules/analytics/delivery`
- [ ] DashboardReadModel
- [ ] TicketTrendReadModel
- [ ] AgentPerformanceReadModel
- [ ] SatisfactionTrendReadModel
- [ ] SLATrendReadModel
- [ ] 先做定时聚合
- [ ] 再做事件驱动增量聚合

验收：

- 统计从业务 service 中分离

## B8 voice

- [ ] 创建 `internal/modules/voice/domain`
- [ ] 创建 `internal/modules/voice/application`
- [ ] 创建 `internal/modules/voice/infra`
- [ ] 创建 `internal/modules/voice/delivery`
- [ ] 定义 `CallSession`
- [ ] 定义 `MediaSession`
- [ ] 定义 `VoiceParticipant`
- [ ] 定义 `Recording`
- [ ] 定义 `Transcript`
- [ ] 实现 StartCall
- [ ] 实现 AnswerCall
- [ ] 实现 EndCall
- [ ] 实现 TransferCall
- [ ] 发布 `call.started`
- [ ] 发布 `call.ended`

验收：

- voice 模块作为 SIP 的业务落点
