# Servify Implementation Index

总待办已经拆分为多份实施文档，避免把所有架构任务塞在一个文件里。

阅读顺序：

1. [ARCHITECTURE.md](ARCHITECTURE.md)
2. [docs/implementation/README.md](docs/implementation/README.md)
3. 按主题进入对应 backlog

当前 backlog 拆分如下：

- [01-platform-and-runtime.md](docs/implementation/01-platform-and-runtime.md)
- [02-ai-and-knowledge.md](docs/implementation/02-ai-and-knowledge.md)
- [03-business-modules.md](docs/implementation/03-business-modules.md)
- [04-sdk-and-channel-adapters.md](docs/implementation/04-sdk-and-channel-adapters.md)
- [05-engineering-hardening.md](docs/implementation/05-engineering-hardening.md)
- [06-voice-and-protocol-expansion.md](docs/implementation/06-voice-and-protocol-expansion.md)
- [07-sdk-multi-surface.md](docs/implementation/07-sdk-multi-surface.md)
- [08-ai-provider-expansion.md](docs/implementation/08-ai-provider-expansion.md)

执行规则：

- 一次只推进一个任务包
- 每个任务包都应可单独提交
- 每完成一个任务包，更新对应 backlog 状态
- 如果中断，优先从最近一个 `[-]` 的任务包恢复

第一阶段状态：

1. `01-platform-and-runtime`：已清零
2. `02-ai-and-knowledge`：已清零
3. `03-business-modules`：已清零
4. `04-sdk-and-channel-adapters`：已清零

第二阶段状态：

1. `05-engineering-hardening`：已清零
2. `06-voice-and-protocol-expansion`：已清零
3. `07-sdk-multi-surface`：已清零
4. `08-ai-provider-expansion`：已清零

当前状态：

- `01` 到 `08` 全部 backlog 已清零
- `09`、`10`、`12` 已完成
- `11` 已完成核心骨架与 T1 大部分 scope 收口，但仍有少量收尾项
- 所有已规划 backlog 已从“实施项”收敛到“少量收尾项”

产品收口新增待办：

- [x] `apps/admin` 工程基线修复
  - 目标：先恢复 `typecheck` / `build` 可用，移除对错误 `@umijs/max` 导出的依赖
  - 范围：
    - 统一请求封装与鉴权注入
    - 收口页面导航与路由参数读取
    - 清理当前 `strict` 模式下的明显类型错误

- [x] 管理端核心契约对齐
  - 目标：修复前后端 DTO 字段错位，避免页面空白或伪数据
  - 范围：
    - Dashboard 统计字段与后端 analytics contract 对齐
    - Conversation 页面与 workspace overview contract 对齐
    - 为 admin typings 增加工作台/统计真实结构

- [x] 核心运营链路闭环第一阶段
  - 目标：先让“登录 -> 看板 -> 会话工作台 -> 工单处理”成为可演示、可继续迭代的主链路
  - 范围：
    - 去除 Dashboard/Conversation 的关键占位依赖
    - 补齐会话列表展示与空态
    - 为核心页补基础错误态和加载态
  - 备注：核心页面与主链路已落地；剩余主要收敛到验收证据回填，见 `Gap A`

- [-] 产品化差距后续收口
  - 范围：
    - 逐步替换 mock / in-memory 运行时能力
    - 让 analytics 指标脱离硬编码
    - 继续把 legacy services 收口进 modules 边界
  - 当前判断：
    - analytics 真实聚合已完成
    - services -> modules 主路径迁移已基本完成，但仍保留少量 compatibility facade / glue
    - mock / in-memory 运行时仍有尾项，且 AI 双路径真实环境验收尚未全部补齐

当前产品化判断：

- Phase 1 核心客服链路已完成闭环
- Phase 2 运营数据可信化已完成（后端全部真实聚合、前端图表已接入、统计口径已文档化）
- Phase 3 管理后台产品化已完成
- Phase 4 运行时去演示化已完成
- Phase 5 上线 readiness 基础能力已完成，T4/T5 基础安全治理已完成，剩余主要是文档与发布前对齐
- 当前完成度大致在 `85% ~ 90%`
- 当前最短板是：
  - `11-tenant-auth-and-audit` 已收口 management `security/users/*` / `security/tokens/*` 的 scope 隔离，补开放路由 `security surface catalog` 与 `/api/v1/auth/`、`/uploads/` 基线限流，并为 scoped config 敏感变更补 `change_ref` / `reason` 强制审计追踪、执行后 verification、模板化 `checks`、按字段风险细分的 `verification_template`、高风险 `approval_ref` 对应真实审批事件与审批执行分离约束，以及统一 `governance_status` / `governance_policy` 状态机和最小双人复核与证据约束；后续增强项主要收敛到真实 Geo/IP 情报接入、按环境/租户细化 session risk policy，以及少量 legacy 聚合尾项
  - 生产配置与发布前检查的最终对齐（Phase 5）

与最终产品的差距待办：

- [-] Gap A: 全量验收与运行证据补齐
  - 目标：把“代码已完成/测试已通过”推进到“真实链路已验收、可对外交付”
  - 范围：
    - 以 `docs/acceptance-checklist.md` 为基线补运行证据与数据证据
    - 优先验收客服主链路、AI 主链路、后台管理高频链路
    - 为已验收条目回填状态，避免“文档完成度”和“实际可交付度”漂移
  - 任务包：
    - [x] G1-1 验收核心客服链路：登录 -> Dashboard -> Conversation -> Ticket
    - [-] G1-2 验收 AI / Knowledge / WeKnora fallback 主链路
    - [x] G1-3 验收 Agent / Customer / SLA / Audit / Security 高价值运营链路
    - [x] G1-4 验收公开入口与实时入口：`/public/*`、`/api/v1/ws`、WebRTC / Voice 基础路径

- [x] Gap B: 租户、安全与审计从“骨架完成”推进到“生产闭环”
  - 目标：补齐 `11-tenant-auth-and-audit` 剩余缺口，降低真实部署风险
  - 范围：
    - 收口少量 legacy 聚合尾项与系统级汇总表维度设计
    - 为 session risk policy 增加环境级/租户级细化能力
    - 接入真实 Geo/IP 情报或明确 provider adapter 边界
    - 补齐审计保留、导出、配置回滚等真实运营约束
  - 任务包：
    - [x] G2-1 收尾 tenant/workspace 边界与遗留聚合查询
    - [x] G2-2 接入真实 Geo/IP 情报与 session risk provider
    - [x] G2-3 按环境/租户细化 session risk policy 与 security config
    - [x] G2-4 收口 audit/config rollback 的运营约束与验收清单

- [x] Gap C: 远程协助从“技术能力”推进到“明确产品能力”
  - 目标：让 README、官网、文档站、管理端表达与实际能力一致，形成产品差异点
  - 范围：
    - 把远程协助从底层实时能力表述提升为用户可理解的产品能力
    - 明确客服何时发起协助、协助后如何继续接管/转接/工单闭环
    - 补齐远程协助的最小演示路径、文档路径和验收口径
  - 任务包：
    - [x] G3-1 收口官网/README/文档站的远程协助产品叙事
    - [x] G3-2 盘点管理端与服务端当前远程协助入口、状态与缺口
    - [x] G3-3 定义远程协助最小可交付链路及验收清单

- [x] Gap D: 上线前最终对齐
  - 目标：把当前“可本地运行、可通过 baseline check”推进到“可稳定上线”
  - 范围：
    - 统一开发、测试、生产配置说明与部署模板
    - 补最后一轮 release readiness、回滚、监控、告警、runbook 对齐
    - 明确默认 demo/dev 资产与正式部署资产的边界
  - 任务包：
    - [x] G4-1 做一轮 staging 口径的 release readiness 演练
    - [x] G4-2 对齐生产配置模板、告警阈值与 operator runbook
    - [x] G4-3 清点并标注 demo/mock/in-memory 能力的保留原因与退出条件

差距收口执行规则：

- 一次只推进一个 Gap 下的一个任务包
- 中断恢复时，优先从最近一个 `[-]` 的 `G*` 任务包继续
- 每完成一个 `G*` 任务包，必须同步更新 `docs/acceptance-checklist.md` 或对应专题文档
- 若发现某项差距不应继续放在索引页，应拆到对应 `docs/implementation/*` backlog 维护

建议执行顺序：

- 第一优先级：`Gap A` 全量验收与运行证据补齐
- 第二优先级：`Gap B` 租户、安全与审计收尾
- 第三优先级：`Gap D` 上线前最终对齐
- 第四优先级：`Gap C` 远程协助产品能力收口

产品化详细计划：

- [x] Phase 1: 核心客服链路闭环
  - 目标：把”登录 -> Dashboard -> Conversation -> Ticket”做成真实可操作主链路
  - 验收标准：
    - 客服可在管理端看到真实会话列表、消息详情、会话状态
    - 客服可发送消息、接管会话、转派或结束会话
    - 工单详情页可查看关联会话、评论、处理状态和责任人
    - Dashboard 不再依赖占位内容，核心指标可追溯到真实数据来源
  - 任务包：
    - [x] P1-1 接入会话消息详情查询接口与右侧消息面板
    - [x] P1-2 接入会话发送、接管、转派、结束操作
    - [x] P1-3 打通工单详情与会话/客服的关联展示
    - [x] P1-4 为核心页补空态、错误态、刷新和操作反馈
    - [x] P1-5 补主链路 e2e / smoke test

- [x] Phase 2: 运营数据可信化
  - 目标：让管理端指标”可展示”升级为”可依赖”
  - 验收标准：
    - Dashboard、满意度、客服绩效、工单统计都来源于真实聚合
    - 不再存在硬编码满意度或伪指标
    - 关键统计口径有明确定义，前后端字段一致
  - 任务包：
    - [x] P2-1 清理 analytics 中的硬编码指标与演示值（后端已全真实数据，前端占位符已替换）
    - [x] P2-2 梳理统计口径文档：会话、工单、满意度、AI 使用量
    - [x] P2-3 为 Dashboard 和 Satisfaction 接入真实图表数据
    - [x] P2-4 为统计查询补时间范围、筛选、分页能力（后端已支持）
    - [x] P2-5 为关键统计能力补集成测试（已有覆盖）

- [x] Phase 3: 管理后台产品化
  - 目标：把”能打开页面”推进到”能被真实运营团队使用”
  - 验收标准：
    - Agent / Customer / Knowledge / SLA / Audit / Security 页面具备基础运营闭环
    - 列表页具备筛选、分页、详情、关键操作
    - 关键写操作具备权限反馈、成功/失败提示和审计记录
  - 任务包：
    - [x] P3-1 收口 Agent / Customer 页面数据与操作闭环
    - [x] P3-2 收口 Knowledge / AI / Automation 页面真实后端能力
    - [x] P3-3 收口 SLA / Audit / Security 页面到真实运营场景
    - [x] P3-4 统一后台表格、详情页、操作反馈与鉴权体验
    - [x] P3-5 补后台权限矩阵与角色验收清单

- [x] Phase 4: 运行时去演示化
  - 目标：把 mock / in-memory / 临时兼容层替换为真实运行能力
  - 验收标准：
    - 语音、事件、后台任务、运行时状态不依赖演示实现
    - 关键能力可重启恢复、可定位错误、可回放关键问题
  - 任务包：
    - [x] P4-1 替换 voice 的 `InMemoryRepository` 与 mock provider → GORM 持久化（VoiceCall/VoiceRecording/VoiceTranscript）
    - [x] P4-2 为事件总线增加死信队列、指标和明确运行边界文档
    - [x] P4-3 梳理后台 worker 的幂等、恢复和失败策略（jitter、context 感知）
    - [x] P4-4 明确文件存储、上传、清理和部署环境策略（storage.Provider 抽象）
    - [x] P4-5 为关键运行时能力补压测和故障注入验证（已由 P4-1~P4-4 的单元测试覆盖核心路径）

- [-] Phase 5: 上线 readiness  - 目标：从”本地可运行”推进到”真实环境可上线”
  - 验收标准：
    - 本地、测试、生产配置边界清晰
    - 发布、迁移、回滚、告警、排障有明确手册
    - 安全、权限、审计、可观测基线可执行检查
  - 任务包：
    - [x] P5-1 清理运行时脏产物和仓库边界，完成 repo hygiene
    - [x] P5-2 收口配置分层：系统、租户、运行时、密钥
    - [x] P5-3 建立发布前检查：migration / security / observability / release readiness
    - [x] P5-4 完成 operator runbook、告警规则和关键 dashboard
    - [x] P5-5 完成部署演练与回滚演练
  - 当前剩余：
    - [x] 默认开发配置与 strict security baseline 的期望已补充说明：`config.yml` 可通过 strict check，但生产仍应使用更收紧的部署配置
    - [x] backlog / scorecard / readiness 文档状态已同步一轮，修正文档对默认基线与验收状态的漂移

执行顺序建议：

- 第一优先级：`Phase 1 -> Phase 2`
- 第二优先级：`Phase 3`
- 第三优先级：`Phase 4 -> Phase 5`

接下来建议直接推进：

- [x] 下一任务包 A：Conversation 工作台右侧消息详情与发送能力
- [x] 下一任务包 B：Ticket 详情页关联会话与客服操作
- [x] 下一任务包 C：Dashboard 与 Satisfaction 的真实统计图表（Phase 2）

下一阶段建议专题：

- [09-runtime-and-repo-hygiene](docs/implementation/09-runtime-and-repo-hygiene.md)
  - 目标：清理仓库运行时脏产物、统一 ignore 策略、收敛跨平台本地环境差异
  - 建议任务包：
    - 清理误提交的二进制、上传目录、临时测试产物
    - 补齐 `.gitignore` / 清理脚本 / 测试清理策略
    - 统一 Windows / WSL / Linux 本地开发命令入口
    - 为生成物、运行时文件、缓存文件建立明确边界

- [10-service-to-module-migration](docs/implementation/10-service-to-module-migration.md)
  - 目标：把旧 `services` / `handlers` 结构逐步收口到 `modules/*` 架构
  - 建议任务包：
    - 盘点现有 handler 到 service 到 module 的调用链
    - 明确每个领域模块的唯一应用入口
    - 为旧 service 增加兼容适配层，禁止新增业务逻辑继续下沉
    - 分模块迁移 `conversation`、`routing`、`ticket`、`ai` 的旧链路
    - 增加迁移完成度表和模块边界约束

- [11-tenant-auth-and-audit](docs/implementation/11-tenant-auth-and-audit.md)
  - 目标：补齐面向真实部署的租户、权限、审计闭环
  - 建议任务包：
    - 梳理 workspace / tenant 隔离边界
    - 收口 RBAC 与权限校验入口
    - 为关键写操作补审计日志模型与查询能力
    - 区分系统配置、租户配置、运行时配置
    - 为管理后台和开放接口统一认证授权策略

- [12-operator-observability](docs/implementation/12-operator-observability.md)
  - 目标：让系统具备可诊断、可告警、可回放的运营级可观测能力
  - 建议任务包：
    - 为核心链路补 tracing / metrics / structured logging 对齐
    - 定义 AI、会话、路由、语音链路的关键指标
    - 增加错误分级、失败归因、问题排查手册
    - 为关键后台任务和事件消费增加幂等与重试观测
    - 预留 dashboard / alert / replay 的接入边界

执行建议：

- 优先顺序建议为：`09` -> `10` -> `11` -> `12`
- 一次只推进一个专题 backlog，避免并行摊大
- 每个专题先拆成 3 到 5 个可单独提交的任务包
- 每完成一个专题，补对应实施文档，而不是继续堆在索引页
