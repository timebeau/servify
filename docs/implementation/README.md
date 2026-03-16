# Servify Implementation Backlogs

本目录按主题拆分架构实施任务。

目的：

- 让每份 backlog 保持可读
- 让 server、AI、SDK、协议接入各自可独立推进
- 支持中断恢复

文件说明：

- [01-platform-and-runtime.md](./01-platform-and-runtime.md)
  - 入口、bootstrap、router、auth、event bus、realtime 等平台任务

- [02-ai-and-knowledge.md](./02-ai-and-knowledge.md)
  - provider、AI orchestration、tooling、knowledge indexing 等任务

- [03-business-modules.md](./03-business-modules.md)
  - conversation、routing、agent、ticket、customer、automation、analytics、voice

- [04-sdk-and-channel-adapters.md](./04-sdk-and-channel-adapters.md)
  - sdk core、web sdk、future api/app sdk 预留、channel adapters、SIP adapter

状态约定：

- `[ ]` 未开始
- `[-]` 进行中
- `[x]` 已完成

执行建议：

1. 先做 `01-platform-and-runtime`
2. 再做 `02-ai-and-knowledge`
3. 再拆 `03-business-modules`
4. 并行规划 `04-sdk-and-channel-adapters`，但当前只实现 web sdk
