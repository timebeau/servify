# Servify Execution Todo

这个文件现在作为“可中断、可恢复”的总控待办，不再只是索引页。

使用规则：

1. 先处理 `P0-代码审查问题`，再处理 `P1-上线与交付闭环`，最后处理 `P2/P3-企业级增强`。
2. 一次只把一个任务推进到“有明确产出”的状态：代码、测试、文档、证据，至少完成一类。
3. 每次中断前必须更新：
   - `状态`
   - `最近进展`
   - `下一步`
   - `阻塞项`
4. 恢复执行时，优先从最近一个 `[-]` 项继续，其次处理最近一个 `[!]` 项。
5. 如果某项已经拆到专题文档，`todo.md` 仍保留摘要、优先级和恢复指针，避免执行上下文丢失。

状态约定：

- `[ ]` 未开始
- `[-]` 进行中 / 中断后从这里恢复
- `[x]` 已完成
- `[!]` 高优先级问题，优先于普通 backlog
- `[?]` 需进一步确认

关联文档：

1. [README.md](README.md)
2. [ARCHITECTURE.md](ARCHITECTURE.md)
3. [docs/acceptance-checklist.md](docs/acceptance-checklist.md)
4. [docs/delivery-priorities.md](docs/delivery-priorities.md)
5. [docs/implementation/README.md](docs/implementation/README.md)

---

## 当前结论

基于当前仓库代码审查，现阶段最优先的问题不是“功能完全没写”，而是以下几类真实交付风险：

1. 运行时仍保留生产路径中的 `InMemory` / `mock` / `legacy` 兼容实现，导致“本地可演示”与“生产可交付”之间仍有断层。
2. 验收清单里仍存在一批主链路 `部分通过 / 未验 / 阻塞` 条目，说明代码存在不等于已可交付。
3. 启动、配置、事件、语音、客服运行态等关键基础设施还有明确的企业级硬伤，应先收口再继续扩展功能面。

---

## P0 代码审查问题

这些问题来自本轮直接审查代码后的判断，应优先进入执行序列。

### [!] P0-1 事件总线仍默认使用进程内 `InMemoryBus`

- 现状：
  - `apps/server/cmd/server/main.go` 直接使用 `eventbus.NewInMemoryBus()`
  - `apps/server/internal/app/bootstrap/app.go` 默认也初始化 `eventbus.NewInMemoryBus()`
- 风险：
  - 事件不持久化，进程重启即丢失 in-flight event
  - 自动化、统计、审计等异步能力的交付边界仍依赖单进程存活
  - 不利于企业部署下的重启恢复、扩缩容、可回放排障
- 代码证据：
  - `apps/server/cmd/server/main.go`
  - `apps/server/internal/app/bootstrap/app.go`
  - `apps/server/internal/platform/eventbus/inmemory_bus.go`
- 执行要求：
  - 明确“生产支持的事件总线边界”
  - 如果短期不引入外部 MQ，至少要把 runtime boundary、失败补偿、恢复策略、运维告警写实
  - 若引入可持久化总线，需要保留兼容接口并补回归测试
- 验收标准：
  - 明确 dev/demo 与 prod 的事件总线策略
  - 至少一条异步链路可证明重启后不会静默丢失关键业务结果，或明确声明当前不承诺 durability 并落实监控/告警/死信审计
- 状态：`[-]`
- 最近进展：已确认启动入口和 bootstrap 默认都走 `InMemoryBus`
- 下一步：梳理当前 event bus 使用面，决定是“先文档化边界并强化告警”还是“直接实现持久化总线适配层”
- 阻塞项：暂无

### [!] P0-2 Voice 运行时仍依赖 mock provider

- 现状：
  - `apps/server/internal/app/server/runtime.go` 中 `RecordingService` / `TranscriptService` 仍注入 `voice/provider/mock`
- 风险：
  - 录音、转写链路目前不具备真实生产语义
  - 虽然仓储已落 GORM，但 provider 仍是 mock，会让“已持久化”掩盖“未真正接入外部能力”
- 代码证据：
  - `apps/server/internal/app/server/runtime.go`
  - `apps/server/internal/modules/voice/provider/mock/*`
- 执行要求：
  - 明确 provider 抽象和真实 provider 的接入方式
  - 至少把 mock provider 从默认 prod runtime 中剥离，改为显式 dev/test 配置
- 验收标准：
  - 默认生产配置下不再隐式使用 mock provider
  - dev/test 配置与 prod 配置边界清晰
  - voice 基础链路至少有一条真实 provider 或明确的 no-op/disabled product boundary
- 状态：`[ ]`
- 最近进展：已确认 runtime 默认 wiring 仍注入 mock provider
- 下一步：检查 voice provider 抽象是否足够支撑真实实现接入
- 阻塞项：真实外部 voice provider 选型尚未确认

### [!] P0-3 Agent 运行态仍依赖内存注册表与 legacy runtime 适配层

- 现状：
  - `apps/server/internal/services/agent_service_assembly.go` 仍使用 `agentinfra.NewInMemoryRegistry()`
  - `legacyRuntime` / `maintenance sync` 仍存在
- 风险：
  - 客服在线状态、负载、会话分配仍偏向单进程运行时模型
  - 重启、扩容、多实例下的一致性边界不够明确
- 代码证据：
  - `apps/server/internal/services/agent_service_assembly.go`
  - `apps/server/internal/services/agent_service.go`
  - `apps/server/internal/services/agent_legacy_runtime_adapter.go`
- 执行要求：
  - 明确 agent presence / load / assignment 的真实来源
  - 决定是持久化 registry、外部协调层，还是显式声明单实例约束
- 验收标准：
  - 客服上下线、负载、会话接管在重启后行为明确
  - 文档与实现对齐，不再让内存态伪装成企业能力
- 状态：`[ ]`
- 最近进展：已确认 assembly 默认 wiring 仍依赖内存 registry
- 下一步：梳理 admin/transfer/session assignment 依赖哪些 runtime 状态
- 阻塞项：多实例策略未定

### [!] P0-4 配置加载仍存在直接 `panic`，且默认模型配置偏旧

- 现状：
  - `apps/server/internal/config/config.go` 的 `Load()` 在 `viper.Unmarshal` 失败时直接 `panic`
  - 默认 OpenAI 模型仍为 `gpt-3.5-turbo`
- 风险：
  - 配置错误时服务启动失败不可控，不利于企业部署排障
  - 默认模型值容易与当前产品策略、真实支持矩阵脱节
- 代码证据：
  - `apps/server/internal/config/config.go`
- 执行要求：
  - 把配置错误改为显式返回错误，由上层决定退出方式
  - 清理默认模型策略，避免误导性默认值进入生产
- 验收标准：
  - 配置解析错误可被测试覆盖、日志可读
  - 默认 AI 配置与 README / docs / 实际 provider 策略一致
- 状态：`[ ]`
- 最近进展：已确认存在 `panic(err)` 路径
- 下一步：梳理 `LoadConfig` 调用链，设计非 panic 返回路径
- 阻塞项：需确认是否有依赖旧签名的调用方

### [!] P0-5 Fallback 配置仍暴露 `legacy_kb_enabled`，兼容语义未完全收口

- 现状：
  - 配置结构中仍有 `Fallback.LegacyKBEnabled`
  - 服务与文档中仍存在大量 `legacy` / `WeKnora compatibility` 语义
- 风险：
  - 对外能力命名与内部实现命名混用
  - 长期会放大配置理解成本和运维误判
- 代码证据：
  - `apps/server/internal/config/config.go`
  - `apps/server/internal/services/ai.go`
  - `apps/server/internal/services/ai_enhanced.go`
  - `apps/server/internal/handlers/ai_handler.go`
- 执行要求：
  - 把“兼容实现”与“对外能力”命名分层
  - 逐步清理公开配置中的 `legacy` 表述
- 验收标准：
  - 面向用户/运维的配置与状态接口不再以 `legacy` 作为核心能力命名
  - compatibility 路径只保留在内部实现或迁移文档中
- 状态：`[ ]`
- 最近进展：已确认 `legacy_kb_enabled` 仍在配置主结构中
- 下一步：盘点哪些 API/文档/测试仍依赖 legacy 命名
- 阻塞项：需兼容已有配置文件

---

## P1 上线与交付闭环

这些是“代码已有，但必须补齐证据或收尾”的上线级事项。

### [!] P1-1 AI / Knowledge 主链路验收闭环

- 目标：
  - 把 `upload` / `sync` / `enable-disable` / `fallback` 从“部分通过”推进到“通过”
- 关联文档：
  - `docs/acceptance-checklist.md`
  - `docs/delivery-priorities.md`
- 关键动作：
  - 补真实 provider 场景下的运行证据
  - 区分 Dify 主路径与 WeKnora compatibility 路径
  - 为失败、超时、fallback 生成可留档证据
- 验收标准：
  - 至少一条真实文档上传、同步、查询命中成功
  - fallback 有实际日志、响应、状态三类证据
- 状态：`[-]`
- 最近进展：现有验收已覆盖本地 provider mock 和部分控制面接口
- 下一步：以 `docs/acceptance-checklist.md` 为准把 G1-2 补成“通过”
- 阻塞项：外部 provider 环境与凭证

### [!] P1-2 Auth 自助 session 链路补齐真实验收

- 范围：
  - `login`
  - `refresh`
  - `sessions`
  - `logout-current`
  - `logout-others`
- 风险：
  - 企业客户会把这组能力视为安全基本面
- 验收标准：
  - refresh token 轮转成功
  - 旧 refresh token 复用失败
  - 当前会话与其它会话退出都能看到状态变化
- 状态：`[ ]`
- 最近进展：自动化存在，但验收矩阵仍为 `部分通过`
- 下一步：补真实请求证据和至少一条反例证据
- 阻塞项：暂无

### [!] P1-3 会话工作台主操作补齐到“通过”

- 范围：
  - 会话详情
  - 发消息
  - 接管
  - 转接
  - 关闭
- 风险：
  - 客服主链路不能长期停留在“整体通过、分项阻塞”
- 验收标准：
  - `GET /api/omni/sessions/:id` 不再阻塞
  - 每个操作有请求结果与状态变化证据
- 状态：`[ ]`
- 最近进展：验收矩阵中 `会话详情` 仍为 `阻塞`
- 下一步：定位 `conversation_workspace_handler` 当前剩余阻塞原因
- 阻塞项：需结合当前 seed 数据与读模型一致性排查

### [!] P1-4 运行基线最小事实补齐

- 范围：
  - `GET /ready`
  - `GET <metricsPath>`
  - `make build`
  - `GET /api/v1/messages/platforms`
- 风险：
  - 这些能力未验会直接影响上线口径
- 验收标准：
  - 全部回填到 `docs/acceptance-checklist.md`
- 状态：`[ ]`
- 最近进展：验收矩阵中仍有多个 `未验`
- 下一步：先跑最小可复现命令并回填证据
- 阻塞项：暂无

### [ ] P1-5 Ticket 主闭环剩余高频操作补齐

- 范围：
  - 更新工单
  - 评论
  - 关闭
  - 统计
  - 导出
- 验收标准：
  - 对应 API 从 `未验` 或 `部分通过` 推进到 `通过`
- 状态：`[ ]`
- 最近进展：工单主链路基础已通，但高频运营动作还未全部验收
- 下一步：按验收矩阵逐项补证据
- 阻塞项：暂无

---

## P2 企业级项目差距

这些不是当前最紧急代码缺陷，但属于“距离企业级项目”的重点 backlog。

### [ ] P2-0 客户侧推荐问题与上下文联想问题

- 当前判断：
  - 已有后台辅助推荐接口 `GET/POST /api/assist/suggest`
  - 当前返回能力主要是 `intent`、`similar_tickets`、`knowledge_docs`
  - 该能力挂在 management `assist` 权限路由下，更像后台辅助检索，不是客户咨询入口的正式产品能力
  - 当前未看到 Web 客户侧 / SDK / demo 已接入“首屏推荐问题”或“基于上下文的下一问联想”
- 代码证据：
  - `apps/server/internal/handlers/suggestion_handler.go`
  - `apps/server/internal/services/suggestion_service.go`
  - `apps/server/internal/app/server/router.go`
  - `sdk/packages/core/src/sdk.ts`
- 产品目标：
  - 客户刚进入咨询页时，能看到可点击的推荐问题
  - 客户发起几轮对话后，系统可基于当前上下文动态联想下一问
  - 推荐问题可与知识库、历史工单、热门问题、AI 意图识别联动
- 建议拆分：
  - `RQ-1` 首屏热门问题推荐
  - `RQ-2` 会话内上下文联想问题
  - `RQ-3` 客户侧推荐接口与权限边界
  - `RQ-4` Web SDK / demo / 官网接入
  - `RQ-5` 埋点、点击率、转化率与验收口径
- 验收标准：
  - 客户未输入前可拿到一组推荐问题
  - 客户输入后可拿到一组基于上下文变化的联想问题
  - 推荐项可点击进入提问，不只是展示静态文案
  - 前后端、SDK、验收文档有统一口径
- 状态：`[ ]`
- 最近进展：已确认当前只有后台辅助推荐接口，没有完整客户侧链路
- 下一步：先决定推荐策略走“规则+知识库”还是“AI 生成+规则兜底”
- 阻塞项：产品策略、推荐来源与客户侧接口边界尚未定稿

### [ ] P2-1 多实例与高可用边界明确化

- 范围：
  - agent online/runtime state
  - WebSocket / WebRTC 会话路由
  - event bus 异步链路
  - worker 幂等和重启恢复
- 目标：
  - 明确 Servify 当前是单实例优先，还是支持多实例协同
- 验收标准：
  - 文档、部署说明、运行时行为一致

### [ ] P2-2 配置治理与环境分层强化

- 范围：
  - dev/staging/prod 配置模板
  - 敏感配置来源
  - 启动前校验
  - 配置漂移检查
- 目标：
  - 防止示例配置、开发默认值和生产配置混用
- 验收标准：
  - 配置加载、校验、模板、文档完全对齐

### [ ] P2-3 数据恢复、备份与迁移演练

- 范围：
  - 数据库迁移回滚
  - 审计/工单/会话关键表恢复
  - 上传文件与知识文档资产恢复
- 目标：
  - 从“可迁移”升级到“可恢复”
- 验收标准：
  - 至少一轮备份恢复演练证据

### [ ] P2-4 可观测性从“有指标”升级到“可运维”

- 范围：
  - 关键业务 SLI/SLO
  - 异步失败告警
  - AI/provider 失败分类
  - 远程协助与实时链路诊断
- 目标：
  - 运维能在故障时快速定位，不依赖人工翻日志
- 验收标准：
  - 告警规则、dashboard、runbook 三者一致

### [ ] P2-5 安全治理继续收口到首批企业交付标准

- 范围：
  - session 风险策略
  - refresh token 治理
  - 审批与回滚证据
  - 公开接口治理
- 目标：
  - 从“最小可部署治理”提升到“企业试点可接受”
- 验收标准：
  - 对应安全面能力均有真实验收，不只靠单测

### [ ] P2-6 管理端产品化收尾

- 范围：
  - Satisfaction
  - Session Transfer
  - Customer/Agent/Workspace 运营细节
  - Ticket/Statistics 细分运营查询
- 目标：
  - 让运营团队可持续日常使用，而不只是跑演示
- 验收标准：
  - 高频后台功能不存在大面积 `未验`

### [ ] P2-7 SDK 与多端 contract 稳定性治理

- 范围：
  - `sdk/packages/*`
  - surface governance
  - 版本同步
  - 示例可运行性
- 目标：
  - 防止 API 面和 SDK 面发生漂移
- 验收标准：
  - SDK smoke test、example smoke test、surface governance 全绿

### [ ] P2-8 性能、容量与压测基线

- 范围：
  - WebSocket 连接数
  - AI 查询延迟
  - ticket/conversation 高并发读写
  - 文件上传与知识同步
- 目标：
  - 从“能跑”提升到“知道能承受多少负载”
- 验收标准：
  - 至少有一版容量基线和压测结论

---

## P3 架构与技术债清理

### [ ] P3-1 清理公开语义中的 legacy/compat 混名

- 目标：
  - 内部兼容层继续保留，但对外命名统一到产品语义

### [ ] P3-2 收拢 services 与 modules 的最终边界

- 目标：
  - 继续减少 glue code 和 facade 长期滞留

### [ ] P3-3 清理 demo/mock/in-memory 资产的默认暴露面

- 目标：
  - 默认生产路径不再误接入 demo/mock 能力

### [ ] P3-4 统一运行时装配方式

- 范围：
  - `cmd/server`
  - `cmd/cli/run`
  - `cmd/cli/run_enhanced`
- 目标：
  - 减少重复 wiring 与配置漂移

---

## 执行顺序

建议严格按这个顺序推进：

1. `P0-1` 事件总线边界收口
2. `P0-2` Voice mock provider 剥离
3. `P0-3` Agent runtime 内存态收口
4. `P0-4` 配置 panic 与默认值治理
5. `P0-5` legacy 配置命名清理
6. `P1-1` AI / Knowledge 验收闭环
7. `P1-2` Auth session 验收闭环
8. `P1-3` 会话工作台收口
9. `P1-4` 基线事实补齐
10. `P1-5` Ticket 高频操作验收
11. 再进入 `P2/P3`

---

## 中断恢复模板

每次执行某个任务前，把对应条目下面更新为：

- 状态：`[-]`
- 最近进展：一句话说明已完成什么
- 下一步：一句话说明下次继续做什么
- 阻塞项：没有就写“暂无”

如果任务完成，更新为：

- 状态：`[x]`
- 完成证据：
  - 代码文件
  - 测试命令
  - 文档或验收回填位置

---

## 当前恢复点

- 当前优先恢复任务：`P0-1 事件总线仍默认使用进程内 InMemoryBus`
- 原因：这是当前运行时企业级边界最明确、最基础的短板，且会影响自动化、统计、审计等多个后续能力判断
- 如果本轮无法推进实现，至少先补：
  - 真实边界文档
  - 默认 prod 策略
  - 风险说明
  - 告警/死信/恢复规则
