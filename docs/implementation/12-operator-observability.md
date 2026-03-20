# 12 Operator Observability

范围：

- tracing
- metrics
- structured logging
- 错误分级
- 告警与回放
- 运营诊断能力

## O1 telemetry-conventions

- [ ] 统一 tracing、metrics、logging 的命名约定
- [ ] 为 HTTP、WebSocket、AI、routing、voice 链路定义公共标签
- [ ] 明确 request id、session id、tenant id、trace id 的透传规则
- [ ] 统一结构化日志字段规范

验收：

- 不同模块输出的观测数据可以关联起来，而不是各写各的

## O2 core-service-level-indicators

- [ ] 定义 API、会话、工单、路由、AI、语音的关键指标
- [ ] 明确成功率、延迟、错误率、积压量等 SLI
- [ ] 为后台任务和事件消费定义吞吐与失败指标
- [ ] 为 SDK 与服务端交互补最小体验指标

验收：

- 系统核心健康度可被量化，而不是靠日志猜测

## O3 error-taxonomy-and-diagnostics

- [ ] 定义错误分级与错误类别
- [ ] 区分用户错误、依赖错误、配置错误、系统错误
- [ ] 为关键模块补统一错误映射与日志策略
- [ ] 为常见故障建立排查手册

验收：

- 出问题时可以快速知道“哪类错误、在哪里、怎么查”

## O4 async-reliability-observability

- [ ] 为 event bus 消费、worker、索引任务、AI fallback 增加幂等观测
- [ ] 记录重试次数、死信、跳过、回退路径
- [ ] 为长耗时任务增加阶段性进度与结果摘要
- [ ] 为异步任务建立最小 replay / rerun 接口预留

验收：

- 异步链路从“黑盒”变为可观察、可重试、可解释

## O5 dashboards-alerts-and-replay

- [ ] 盘点需要 dashboard 的核心视角
- [ ] 定义关键告警触发条件与阈值策略
- [ ] 预留会话、AI 请求、路由决策的 replay 边界
- [ ] 为运营排障定义最小诊断面板需求

验收：

- 运营和研发都具备足够的现场信息，不必直接翻原始日志救火
