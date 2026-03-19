# 08 AI Provider Expansion

范围：

- LLM provider 扩展
- Knowledge provider 扩展
- 编排链路稳定性
- AI 可观测性与安全控制

## A1 provider-matrix

- [ ] 盘点当前 `LLMProvider` 与 `KnowledgeProvider` contract 缺口
- [ ] 定义多 provider 配置矩阵
- [ ] 定义优先级、fallback、熔断策略矩阵
- [ ] 为 provider capability 建立声明模型

验收：

- 新增 provider 时不需要靠运行时 if/else 拼接

## A2 llm-provider-expansion

- [ ] 增加第二个真实 `LLMProvider` 适配器占位
- [ ] 统一流式/非流式响应 contract
- [ ] 统一 tool calling / function calling contract
- [ ] 增加 provider 级超时、重试、错误分类

验收：

- LLM provider 不再默认只有 OpenAI 一条路径

## A3 knowledge-provider-expansion

- [ ] 增加第二个真实 `KnowledgeProvider` 适配器占位
- [ ] 抽象索引、删除、重建、检索一致性语义
- [ ] 定义 tenant / knowledge base 多租户映射规则
- [ ] 增加 provider 切换回归测试

验收：

- WeKnora 成为可替换实现，而不是唯一实现

## A4 orchestration-hardening

- [ ] 补 query orchestration 场景测试矩阵
- [ ] 补 prompt builder contract 测试矩阵
- [ ] 补 retrieval 降级路径测试
- [ ] 补 AI fallback 行为文档

验收：

- 编排层行为对 provider 切换保持稳定

## A5 ai-observability-and-policy

- [ ] 增加 provider 维度 tracing/span 标签
- [ ] 增加 token / latency / error 分类指标
- [ ] 预留内容安全与策略拒答 hook
- [ ] 预留 prompt/version 审计记录接口

验收：

- AI 能力具备上线前需要的基本可观测性与策略控制点
