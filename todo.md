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
- `09` 到 `12` 全部 backlog 已清零
- 所有已规划 backlog 均已完成

产品收口新增待办：

- [ ] `apps/admin` 工程基线修复
  - 目标：先恢复 `typecheck` / `build` 可用，移除对错误 `@umijs/max` 导出的依赖
  - 范围：
    - 统一请求封装与鉴权注入
    - 收口页面导航与路由参数读取
    - 清理当前 `strict` 模式下的明显类型错误

- [ ] 管理端核心契约对齐
  - 目标：修复前后端 DTO 字段错位，避免页面空白或伪数据
  - 范围：
    - Dashboard 统计字段与后端 analytics contract 对齐
    - Conversation 页面与 workspace overview contract 对齐
    - 为 admin typings 增加工作台/统计真实结构

- [ ] 核心运营链路闭环第一阶段
  - 目标：先让“登录 -> 看板 -> 会话工作台 -> 工单处理”成为可演示、可继续迭代的主链路
  - 范围：
    - 去除 Dashboard/Conversation 的关键占位依赖
    - 补齐会话列表展示与空态
    - 为核心页补基础错误态和加载态

- [ ] 产品化差距后续收口
  - 范围：
    - 逐步替换 mock / in-memory 运行时能力
    - 让 analytics 指标脱离硬编码
    - 继续把 legacy services 收口进 modules 边界

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
