# 02 AI And Knowledge

范围：

- llm provider
- knowledge provider
- ai orchestrator
- prompt builder
- retriever
- tools
- metrics
- knowledge indexing

## A1 provider contracts

- [x] 创建 `internal/platform/llm`
- [x] 定义 `ChatRequest`
- [x] 定义 `ChatResponse`
- [x] 定义 `TokenUsage`
- [x] 定义 `LLMProvider`
- [x] 创建 `internal/platform/knowledgeprovider`
- [x] 定义 `SearchRequest`
- [x] 定义 `KnowledgeHit`
- [x] 定义 `KnowledgeDocument`
- [x] 定义 `KnowledgeProvider`

验收：

- provider 接口不暴露 OpenAI/WeKnora 私有 DTO

## A2 provider adapters

- [x] 新增 `llm/openai/provider.go`
- [x] 迁移模型 HTTP 调用逻辑
- [x] 新增 `llm/mock/provider.go`
- [x] 新增 `knowledgeprovider/weknora/provider.go`
- [x] 包装 `pkg/weknora`
- [x] 新增 `knowledgeprovider/mock/provider.go`

验收：

- AI 模块可在 mock provider 下运行

## A3 ai module skeleton

- [x] 创建 `internal/modules/ai/domain`
- [x] 创建 `internal/modules/ai/application`
- [x] 创建 `internal/modules/ai/infra`
- [x] 创建 `internal/modules/ai/delivery`
- [x] 定义 `AIRequest`
- [x] 定义 `AIResponse`
- [x] 定义 `TaskType`
- [x] 定义 `RetrievalPolicy`
- [x] 定义 `ToolPolicy`

验收：

- AI 模块独立存在

## A4 query orchestrator

- [x] 新增 `query_orchestrator.go`
- [x] 注入 `LLMProvider`
- [x] 注入 `KnowledgeProvider`
- [x] 实现 happy path
- [x] 实现无知识检索 fallback
- [x] 增加 orchestrator 单测

验收：

- 可以替代旧增强 AI 主流程

## A5 prompt and retrieval

- [ ] 新增 `prompt_builder.go`
- [ ] 拆 system prompt
- [ ] 拆 context prompt
- [ ] 拆 knowledge prompt
- [ ] 新增 `retriever.go`
- [ ] 支持 topK
- [ ] 支持 threshold
- [ ] 支持 strategy

验收：

- prompt 和 retrieval 从旧 service 中独立出来

## A6 tools

- [ ] 定义 `Tool`
- [ ] 定义 `ToolRegistry`
- [ ] 定义 `ToolExecutor`
- [ ] 第一批：
  - [ ] ticket lookup tool
  - [ ] customer lookup tool
  - [ ] handoff tool

验收：

- 工具调用有 schema 和权限检查

## A7 ai guardrails and metrics

- [ ] 输入长度限制
- [ ] 敏感内容拦截占位
- [ ] 输出截断
- [ ] query_count
- [ ] success_count
- [ ] provider_usage_count
- [ ] token_usage
- [ ] latency
- [ ] fallback_count

验收：

- 指标由 orchestrator 统一上报

## A8 knowledge module

- [ ] 创建 `internal/modules/knowledge/domain`
- [ ] 创建 `internal/modules/knowledge/application`
- [ ] 创建 `internal/modules/knowledge/infra`
- [ ] 创建 `internal/modules/knowledge/delivery`
- [ ] 定义 `Document`
- [ ] 定义 `DocumentRepository`
- [ ] 定义 `IndexJobRepository`
- [ ] 实现 CreateDocument
- [ ] 实现 UpdateDocument
- [ ] 实现 DeleteDocument
- [ ] 实现 QueueIndexJob
- [ ] 实现 RunIndexJob

验收：

- 知识库文档和索引成为正式模块
