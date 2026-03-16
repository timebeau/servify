# 04 SDK And Channel Adapters

范围：

- sdk core contracts
- current web sdk refactor
- future api sdk reservation
- future app sdk reservation
- channel adapters
- sip adapter

约束：

- 当前只实现 web sdk
- 但目录和 contracts 必须支持未来 api sdk 和 app sdk

## S1 sdk target structure

- [x] 评审当前 `sdk/packages` 结构
- [ ] 确认现有包映射：
  - [x] `core`
  - [x] `react`
  - [x] `vue`
  - [x] `vanilla`
- [ ] 设计目标结构：
  - [x] `core`
  - [x] `transport-http`
  - [x] `transport-websocket`
  - [x] `web-vanilla`
  - [x] `web-react`
  - [x] `web-vue`
  - [x] `api-client`
  - [x] `app-core`

验收：

- 目标结构文档化
- 不要求一次性重命名现有包

## S2 sdk core contracts

- [ ] 在 `sdk/packages/core` 中增加 contracts 分层
- [ ] 定义 `ClientSession`
- [ ] 定义 `Transport`
- [ ] 定义 `EventEmitter`
- [ ] 定义 `AuthProvider`
- [ ] 定义 `Capability`
- [ ] 定义统一错误模型

验收：

- core 不含浏览器专属 UI 逻辑

## S3 transport split

- [ ] 把 websocket 逻辑从泛化 core 中进一步抽到 transport contract
- [ ] 预留 `transport-http`
- [ ] 预留 reconnect policy
- [ ] 预留 token refresh hook

验收：

- transport 与 framework binding 解耦

## S4 current web sdk refactor

- [ ] `react` 仅依赖 core contracts
- [ ] `vue` 仅依赖 core contracts
- [ ] `vanilla` 仅依赖 core contracts
- [ ] 将浏览器专属逻辑收口到 web binding

验收：

- core 仍可被未来 api/app sdk 复用

## S5 future api sdk reservation

- [ ] 创建 `sdk/packages/api-client` 占位
- [ ] 新增 README 说明目标用途
- [ ] 约定使用场景：
  - [ ] server-to-server api calls
  - [ ] admin automation clients
  - [ ] bot integrations

验收：

- 明确这是保留设计，不实现具体功能

## S6 future app sdk reservation

- [ ] 创建 `sdk/packages/app-core` 占位
- [ ] 新增 README 说明目标用途
- [ ] 约定使用场景：
  - [ ] mobile app session
  - [ ] push token integration
  - [ ] reconnect and offline queue

验收：

- 明确这是保留设计，不实现具体功能

## S7 sdk capability negotiation

- [ ] 定义 capability model
- [ ] 第一批 capability：
  - [ ] chat
  - [ ] realtime
  - [ ] knowledge
  - [ ] remote_assist
  - [ ] voice
- [ ] 在 web sdk 中先实现 capability exposure

验收：

- future api/app sdk 可共享 capability 协议

## S8 channel adapter contracts

- [ ] 在服务端定义统一 `InboundEvent`
- [ ] 在服务端定义统一 `OutboundEvent`
- [ ] 约定 adapter responsibilities
- [ ] web adapter 映射到 conversation
- [ ] 预留 telegram/wecom/whatsapp adapter contract

验收：

- 新接入渠道不需要改 ticket 业务核心

## S9 sip adapter

- [ ] 创建 `internal/platform/sip`
- [ ] 定义 `SIPAdapter`
- [ ] 定义 `InboundCall`
- [ ] 映射 invite -> inbound event
- [ ] 映射 hangup -> inbound event
- [ ] 映射 dtmf -> inbound event

验收：

- SIP 通过 adapter 接入，不进入聊天逻辑分支
