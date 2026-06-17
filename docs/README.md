---
title: 文档首页说明
---

> VitePress 首页入口已经迁移到 [index.md](./index.md)。
>
> 这里保留原始内容入口，避免仓库内引用 `docs/README.md` 失效。

## 产品概览

Servify 不是先做平台再兼容客服，而是先把客服主链路做完整：

`Web 接入 -> AI 首答 -> 远程协助 -> 人工接管 -> 转接协作 -> 工单闭环`

当前文档站的重点，也应该围绕这条链路来理解产品，而不是先从 implementation backlog 开始读。

这里的远程协助，不是把“WebRTC 可用”单独当成功能名词，而是把实时指导、联合排查、人工接管和后续工单处理放进同一条服务链路里。

## 为什么是 Servify

### 独立部署优先

Servify 适合需要独立部署、可控接入和持续演化的企业客服场景。它默认服务的是单个企业团队，而不是把多租户平台化能力放在产品正中央。

### 远程协助是差异点

很多客服系统只做到聊天和工单，但真实服务过程里，经常需要“看着客户做、带着客户做、一起排查问题”。Servify 把远程协助视为产品能力，而不是隐藏在底层实时通信里的附属特性。

当前仓库已经具备这类能力的实时基础，包括 Web 会话、消息、WebSocket / WebRTC 相关链路、人工接管和后续工单闭环；但文档表达会刻意避免把它夸大成“已完整交付某种重型远控产品”。

### AI 和人工不是割裂关系

AI 负责首答、澄清、知识召回和建议，人工可以随时接管、继续处理和协作，让整个客服体验保持连续。

## 推荐阅读

1. [总体架构设计](/ARCHITECTURE)
2. [当前架构分析](/current-architecture)
3. [架构重设计计划](/architecture-redesign-plan)
4. [部署指南](/deployment)
5. [当前交付优先级](/delivery-priorities)
6. [实施计划索引](/implementation/)

## 你可以从这里继续

### 快速了解产品

- [总体架构设计](/ARCHITECTURE)
- [当前架构分析](/current-architecture)
- [架构重设计计划](/architecture-redesign-plan)
- [远程协助](/remote-assistance)
- [部署指南](/deployment)
- [本地开发](/local-development)

### 深入运行与安全

- [运行安全基线](/security-baseline-operations)
- [配置作用域规则](/configuration-scopes)
- [Token 生命周期与密钥轮换](/token-lifecycle-and-key-rotation)
- [开放接口安全清单](/public-surface-security-checklist)
- [当前交付优先级](/delivery-priorities)

### 继续看研发与实施细节

- [实施 Backlog 索引](/implementation/)
- [模块迁移计划](/implementation/10-service-to-module-migration)
- [模块迁移完成度](/implementation/10-migration-scorecard)
- [外部知识库集成指南：Dify 优先，WeKnora 兼容](/WEKNORA_INTEGRATION)
- [v0.1.0 Release Notes](/release-notes-v0.1.0)
- [版本发布策略](/release-versioning)
- [测试金字塔](/testing-pyramid)
- [CI / GitHub Hosted Runner](/CI_SELF_HOSTED)
- [Mermaid 兼容性](/MERMAID_COMPATIBILITY)

## 常用入口

- 本地自检：`make local-check`
- 安全检查：`make security-check CONFIG=./config.yml`
- 可观测性检查：`make observability-check CONFIG=./config.yml`
- 发布检查：`make release-check CONFIG=./config.yml`
