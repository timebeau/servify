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

## A4.1 compatibility bridge

- [x] 新增 legacy `AIServiceInterface` 兼容适配器
- [x] 让旧接口可调用新 `QueryOrchestrator`
- [x] 在标准模式入口 wiring 中逐步替换旧 AI service
- [ ] 在增强模式入口 wiring 中逐步替换旧 AI service

验收：

- 不改现有 handler 对外接口的情况下，可以复用新 AI 模块

## A5 prompt and retrieval

- [x] 新增 `prompt_builder.go`
- [x] 拆 system prompt
- [ ] 拆 context prompt
- [x] 拆 knowledge prompt
- [x] 新增 `retriever.go`
- [x] 支持 topK
- [x] 支持 threshold
- [x] 支持 strategy

验收：

- prompt 和 retrieval 从旧 service 中独立出来

## A6 tools

- [x] 定义 `Tool`
- [x] 定义 `ToolRegistry`
- [x] 定义 `ToolExecutor`
- [ ] 第一批：
  - [x] ticket lookup tool
  - [x] customer lookup tool
  - [x] handoff tool

验收：

- 工具调用有 schema 和权限检查

## A7 ai guardrails and metrics

- [x] 输入长度限制
- [x] 敏感内容拦截占位
- [x] 输出截断
- [x] query_count
- [x] success_count
- [x] provider_usage_count
- [x] token_usage
- [x] latency
- [x] fallback_count

验收：

- 指标由 orchestrator 统一上报

## A8 knowledge module

- [x] 创建 `internal/modules/knowledge/domain`
- [x] 创建 `internal/modules/knowledge/application`
- [x] 创建 `internal/modules/knowledge/infra`
- [x] 创建 `internal/modules/knowledge/delivery`
- [x] 定义 `Document`
- [x] 定义 `DocumentRepository`
- [x] 定义 `IndexJobRepository`
- [x] 实现 CreateDocument
- [x] 实现 UpdateDocument
- [x] 实现 DeleteDocument
- [x] 实现 QueueIndexJob
- [x] 实现 RunIndexJob

验收：

- 知识库文档和索引成为正式模块
