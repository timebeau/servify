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

执行规则：

- 一次只推进一个任务包
- 每个任务包都应可单独提交
- 每完成一个任务包，更新对应 backlog 状态
- 如果中断，优先从最近一个 `[-]` 的任务包恢复

当前推荐开工顺序：

1. `01-platform-and-runtime` 的 `P1 bootstrap skeleton`
2. `01-platform-and-runtime` 的 `P3 event bus`
3. `02-ai-and-knowledge` 的 `A1 provider contracts`
4. `04-sdk-and-channel-adapters` 的 `S1 sdk target structure`
5. `03-business-modules` 的 `B4 ticket query split`
