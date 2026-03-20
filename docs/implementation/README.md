# Servify Implementation Backlogs

本目录按主题拆分架构实施任务。

目的：

- 让每份 backlog 保持可读
- 让 server、AI、SDK、协议接入各自可独立推进
- 支持中断恢复

文件说明：

- [总体架构设计](../ARCHITECTURE.md)
  - 宏观边界、运行时分层、未来扩展方向
- [01-platform-and-runtime.md](./01-platform-and-runtime.md)
  - 入口、bootstrap、router、auth、event bus、realtime 等平台任务

- [02-ai-and-knowledge.md](./02-ai-and-knowledge.md)
  - provider、AI orchestration、tooling、knowledge indexing 等任务

- [03-business-modules.md](./03-business-modules.md)
  - conversation、routing、agent、ticket、customer、automation、analytics、voice

- [04-sdk-and-channel-adapters.md](./04-sdk-and-channel-adapters.md)
  - sdk core、web sdk、future api/app sdk 预留、channel adapters、SIP adapter
- [05-engineering-hardening.md](./05-engineering-hardening.md)
  - CI、测试金字塔、版本发布、文档站点
- [06-voice-and-protocol-expansion.md](./06-voice-and-protocol-expansion.md)
  - voice 协议入口深化、provider 落地、更多常见语音协议预留
- [07-sdk-multi-surface.md](./07-sdk-multi-surface.md)
  - web sdk 收口、future api/app sdk contract 深化、transport 演进
- [08-ai-provider-expansion.md](./08-ai-provider-expansion.md)
  - LLM/knowledge provider 扩展、编排层稳定性、AI 可观测性

状态约定：

- `[ ]` 未开始
- `[-]` 进行中
- `[x]` 已完成

执行建议：

1. 先做 `01-platform-and-runtime`
2. 再做 `02-ai-and-knowledge`
3. 再拆 `03-business-modules`
4. 并行规划 `04-sdk-and-channel-adapters`，但当前只实现 web sdk
5. 第一阶段清零后，进入 `05` 到 `08` 的工程化与扩展阶段

配套专题：

- [版本发布策略](../release-versioning.md)
- [测试金字塔](../testing-pyramid.md)
- [Mermaid 兼容性](../MERMAID_COMPATIBILITY.md)
