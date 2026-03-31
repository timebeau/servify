---
home: true
title: Servify Docs
heroText: Servify Docs
tagline: 架构、实施 backlog、工程化约束与集成指南
actions:
  - text: 查看架构
    link: /ARCHITECTURE
    type: primary
  - text: 进入实施计划
    link: /implementation/
    type: secondary
features:
  - title: 架构与模块边界
    details: 先看整体架构，再进入 platform、AI、业务模块和 voice 扩展任务。
  - title: 工程化约束
    details: CI 质量门禁、生成物一致性、版本发布链路、测试金字塔都在文档站统一维护。
  - title: 集成专题
    details: WeKnora、CI Runner、Mermaid 兼容性和发布策略作为长期维护专题单独收口。
footer: Servify documentation site
---

## 推荐阅读顺序

1. [总体架构设计](/ARCHITECTURE)
2. [实施 Backlog 索引](/implementation/)
3. [工程化硬化与交付](/implementation/05-engineering-hardening)

## 专题文档

- [WeKnora 集成](/WEKNORA_INTEGRATION)
- [配置作用域规则](/configuration-scopes)
- [运行安全基线](/security-baseline-operations)
- [Operator 可观测性](/implementation/12-operator-observability)
- [Token 生命周期与密钥轮换](/token-lifecycle-and-key-rotation)
- [开放接口安全清单](/public-surface-security-checklist)
- [CI / GitHub Hosted Runner](/CI_SELF_HOSTED)
- [版本发布策略](/release-versioning)
- [测试金字塔](/testing-pyramid)
- [Mermaid 兼容性](/MERMAID_COMPATIBILITY)

## 常用校验入口

- `make local-check`
- `make security-check CONFIG=./config.yml`
- `make observability-check CONFIG=./config.yml`
- `make release-check CONFIG=./config.yml`
