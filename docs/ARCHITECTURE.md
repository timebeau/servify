# 架构总览

Servify 当前阶段保持为模块化单体，围绕统一 runtime、AI 编排、业务模块、voice 协议入口和 SDK surface 逐步硬化。

## 阅读入口

- 根文档：仓库根目录 `ARCHITECTURE.md`
- 实施 backlog：[./implementation/README.md](./implementation/README.md)
- 工程化硬化：[./implementation/05-engineering-hardening.md](./implementation/05-engineering-hardening.md)

## 当前文档站中的对应专题

- [实施计划](/implementation/)
- [版本发布策略](/release-versioning)
- [测试金字塔](/testing-pyramid)
- [WeKnora 集成](/WEKNORA_INTEGRATION)

## 站内约定

- 架构总览页只承担导航和阅读顺序，不复制根目录长文
- 详细设计仍以仓库根目录的 `ARCHITECTURE.md` 为主
- 文档站页面负责把实施、工程化和专题文档组织成可部署站点
