# Servify v0.1.0 真实状态审计

> 最后更新: 2026-04-18
> 这份文件不再把“测试通过”直接等同于“发版完成”。

## 当前结论

`v0.1.0` 现在不适合发布。

原因不是“核心功能完全没写”，而是以下三类问题同时存在：

1. `TASKS.md`、`todo.md`、`docs/acceptance-checklist.md`、`docs/release-notes-v0.1.0.md` 的完成口径不一致。
2. 仓库里有一批“接口存在 + 测试通过”，但离“真实发布闭环”还有差距的能力。
3. `0.1.0` blocker 里仍有 `部分通过`、`未验` 或实现语义不够稳定的项。

## 目前可以确认的最小事实

这些能力可以保留“已具备基础能力”表述：

- `make build`、`make release-check CONFIG=./config.yml`、`GET /health`、`GET /ready`、`GET /metrics` 已有明确证据。
- 核心客服主链路已具备可演示性：工作台、会话详情、消息收发、指派、转接、关闭至少已有自动化与部分人工运行证据。
- 工单主闭环具备基础可用性：创建、列表、详情、更新、指派、评论、关闭、统计、导出都已有实现和自动化覆盖。
- Auth session 基础链路不是空壳：登录、refresh、sessions、logout-current、logout-others 已有实现和自动化覆盖。
- 实时基础能力已具备：WebSocket、WebRTC stats / connections、平台消息路由统计有入口和验证。

## 当前不能再写成“已完成”的部分

这些项不是“没有代码”，但不能再按“完成验收”描述：

- AI / Knowledge
  - `upload`、`sync` 在验收矩阵里仍是 `部分通过`
  - fallback 仍主要依赖内存态知识库与规则回复
  - provider 抽象存在未实现空壳与成功语义过宽的问题
  - `knowledge-docs` 已开始持久化 `provider_id/external_id`，`Create/Update` 会同步当前外部 knowledge provider，`Delete` 也会优先使用外部 `document_id`
  - Dify 删除路径已恢复为真实能力，当前能基于持久化的外部 `document_id` 精准删除对应 dataset 文档
  - WeKnora 删除能力仍未闭环，现阶段依然只能按“能力不支持则显式失败并保留本地文档”的保守策略处理
- 后台基础运营面
  - customer、macro、custom-field、statistics、satisfaction、automation 等模块大多只有基础 CRUD / 查询面
  - handler 普遍较薄，错误分类粗，不适合按“高可信完成”描述
- 会话转接运营面
  - 验收矩阵中仍为 `未验`
  - 队列处理与取消等链路存在“部分失败仍返回成功”的风险
- 满意度后台全量运营面
  - 验收矩阵中仍为 `未验`
  - 错误处理仍依赖字符串匹配，稳定性一般
- 宏 / 自动化 / 排班 / 自定义字段
  - 这些能力可以继续留在 backlog，不应再在这个文件里写成“全部完成”

## 当前 blocker

按 `docs/release-0.1.0-acceptance.md` 与 `todo.md`，当前至少还有以下 blocker 未闭环：

1. AI / Knowledge 至少 1 条真实 provider 主路径达到 `通过`
2. Auth 自助 session 链路补齐真实发布证据
3. 会话工作台主操作补齐到 `通过`
4. 运行基线最小事实全部回填
5. Ticket 主闭环高频操作按发布口径补齐证据

## 不可信的旧结论

以下旧结论已不再成立：

- “已完成 37 项”
- “待完成 0 项”
- “测试通过即可视为功能完成”

这些说法会误导后续发布判断，现废止。

## 后续维护规则

1. 以 `todo.md` 和 `docs/acceptance-checklist.md` 作为真实状态源。
2. 只有同时具备代码、自动化、运行、数据四类证据的功能，才允许标记为 `通过`。
3. `TASKS.md` 不再维护“全部完成”式总表，只记录当前真实发布判断。

## 参考验证命令

```bash
go test -count=1 ./apps/server/...
go test -count=1 -tags=integration ./apps/server/internal/handlers ./apps/server/internal/services
make build
make release-check CONFIG=./config.yml
```
