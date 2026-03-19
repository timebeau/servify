# 07 SDK Multi Surface

范围：

- Web SDK 收口
- future API SDK 预留继续细化
- future App SDK 预留继续细化
- 对外契约稳定性

约束：

- 当前只真正实现 web 端
- API/App 仍以 contract、目录、测试替身、说明文档为主

## M1 web-sdk-completion

- [ ] 补 `core` capability negotiation 测试矩阵
- [ ] 补 reconnect policy 实现与测试
- [ ] 补 token refresh hook 实现与测试
- [ ] 补 `react/vue/vanilla` 对统一 contract 的 smoke tests
- [ ] 收敛当前 lint warning

验收：

- Web SDK 达到可持续演进状态，不再依赖隐性行为

## M2 api-sdk-reservation

- [ ] 在 `api-client` 中定义 server-side auth contract
- [ ] 定义 retry/backoff contract
- [ ] 定义 idempotency/request middleware contract
- [ ] 定义 bot/admin automation 使用示例

验收：

- 后续做 server-to-server SDK 时可以直接沿用已有 contract

## M3 app-sdk-reservation

- [ ] 在 `app-core` 中定义 offline queue contract
- [ ] 定义 push token registration contract
- [ ] 定义 reconnect/session restore contract
- [ ] 定义 mobile storage abstraction contract

验收：

- 后续移动端 SDK 不会直接复制 Web SDK 结构

## M4 surface-governance

- [ ] 定义 package 命名与发布策略
- [ ] 定义 public API 审核边界
- [ ] 增加 breaking change checklist
- [ ] 增加 examples 与 package README 对齐检查

验收：

- SDK 对外 surface 有明确治理规则

## M5 transport-evolution

- [ ] 细化 `transport-http` 包结构
- [ ] 细化 `transport-websocket` 包结构
- [ ] 预留 SSE / webhook callback transport contract
- [ ] 预留 API/App/Web 共享 serializer contract

验收：

- transport 层成为真正独立能力，而不是 core 里的零散实现
