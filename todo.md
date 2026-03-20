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
- 后续新增任务应新开专题 backlog，避免回填旧任务包
