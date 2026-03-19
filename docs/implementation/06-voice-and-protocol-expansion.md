# 06 Voice And Protocol Expansion

范围：

- voice 协议统一入口深化
- signaling/media provider 落地
- 录音/转写外部 provider
- 常见语音协议扩展预留

约束：

- `voice` 业务层不能反向依赖具体协议 DTO
- signaling 与 media 必须继续分层

## V1 protocol-end-to-end

- [ ] 为 `sip` 协议入口补 invite -> answer -> hangup 端到端测试
- [ ] 为 `sip-ws` 协议入口补 invite -> dtmf -> hangup 端到端测试
- [ ] 为 `pstn-provider` webhook 入口补 invite -> transfer -> hangup 测试
- [ ] 为 `/api/voice/protocols` 增加契约测试

验收：

- 新协议入口不是只有 adapter 单测，而是有完整 HTTP/runtime/voice 链路验证

## V2 call-control-semantics

- [ ] 在统一协议入口补 hold 事件映射
- [ ] 在统一协议入口补 resume 事件映射
- [ ] 在统一协议入口补 transfer 事件映射到 `voice` use case
- [ ] 明确 DTMF 在 voice/conversation/automation 间的分发策略

验收：

- 常见呼叫控制事件都能走统一归一化模型

## V3 provider-implementation

- [ ] 将录音 provider 从 in-memory 抽成独立 mock/provider 包
- [ ] 将转写 provider 从 in-memory 抽成独立 mock/provider 包
- [ ] 为 provider 增加失败重试和错误模型
- [ ] 预留异步 webhook/callback 回写接口

验收：

- 录音与转写 provider 可以独立替换，不影响 voice 用例

## V4 media-topology

- [ ] 抽象 RTP/SRTP adapter 占位
- [ ] 明确 WebRTC 与 RTP bridge 的 contract
- [ ] 预留 conference/mixer 能力接口
- [ ] 预留 QoS / voice quality metrics 采集点

验收：

- 后续从点对点通话扩展到更复杂媒体拓扑时不需要推翻现有结构

## V5 protocol-reservation

- [ ] 预留 `voiceprotocol.ProtocolH323`
- [ ] 预留 `voiceprotocol.ProtocolFreeSWITCHEventSocket` 或同类 PBX 适配能力
- [ ] 预留 hosted voice vendor webhook adapter contract
- [ ] 为新协议接入编写最小实现模板文档

验收：

- voice 模块不是“只支持 SIP”的隐性设计
