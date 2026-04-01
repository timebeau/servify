# Audit Log Policy

本文件定义当前审计日志的最小查询与保留策略。

## 查询入口

- 管理面查询：`GET /api/audit/logs`

支持过滤：

- `action`
- `resource_type`
- `resource_id`
- `principal_kind`
- `actor_user_id`
- `success=true|false`
- `from=<RFC3339>`
- `to=<RFC3339>`
- `page`
- `page_size`
- 以及通过统一 request scope 投影的 `tenant_id` / `workspace_id`

## 当前记录范围

当前统一覆盖：

- 管理面 `/api/*` 的成功写请求
- 记录主体类型、操作者、动作、资源、请求元数据、请求载荷

当前未完全覆盖：

- 匿名 public surface
- 仅只读查询请求
- 每个 handler 的精细 `before/after` 业务快照
- 审计查询导出

## 保留策略

当前建议基线：

- 默认保留 180 天在线审计日志
- 超过 180 天的记录应归档或清理
- 涉及安全事件、权限变更、关键配置变更的日志应支持更长保留期

当前仓库状态：

- 已具备持久化模型与查询接口
- 已接入定时清理 worker，默认按 180 天保留窗口删除过期记录
- 尚未接入冷热分层或归档存储

## 后续实现建议

1. 为高价值审计事件补冷热分层或外部归档存储，而不是只做删除
2. 为权限变更、配置变更补专门的 `before/after` 快照
3. 为审计查询增加导出能力
