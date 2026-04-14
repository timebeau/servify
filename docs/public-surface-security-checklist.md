# Public Surface Security Checklist

本文用于推进 `11 / T5 security-baseline-for-operations` 中“对外开放接口基础安全清单”。

适用范围：

- `/public/*`
- 匿名 websocket / realtime 建连入口
- 匿名认证入口，如 `/api/v1/auth/login`、`/api/v1/auth/register`、`/api/v1/auth/refresh`
- 公开文件访问入口，如 `/uploads/*`
- 未来开放给第三方或终端用户的公开 API

## 新增开放接口前必须确认

1. 这个接口是否真的需要匿名访问？
2. 是否可以改成 `management surface` 或 `service surface`？
3. 是否会泄漏 tenant、workspace、agent、customer 等内部标识？
4. 是否需要分页、过滤或结果裁剪，避免批量枚举？
5. 是否需要单独路径级 rate limiting，而不是只依赖全局限流？
6. 是否会返回 provider、基础设施或内部拓扑信息？
7. 是否需要审计访问事件或异常流量？

## 最小检查项

### 输入

- 只接受明确 schema 的参数
- 对 query、path、header 做长度和格式限制
- 不信任客户端传入的 tenant / workspace scope

### 输出

- 不返回 secret、token、内部配置、provider 凭据
- 不暴露不必要的内部 ID、数据库主键、文件系统路径
- 错误信息不回显内部堆栈、SQL 或第三方 provider 细节

### 流量控制

- 必须评估是否需要单独路径级 rate limiting
- 对可触发高成本计算、检索、上传的接口使用更严格配额
- 若支持 websocket / streaming，应限制并发连接或单位时间连接数

### 浏览器暴露面

- 检查 CORS 是否只允许必要 origin
- 若接口可被页面嵌入或脚本调用，确认是否需要额外 origin / referer 限制

### 可观测性

- 为异常流量、`429`、`401`、`403` 建立基础监控
- 能通过 `X-Request-ID` 或等效请求标识追踪问题

## 当前接口建议

### `/public/portal/config`

- 仅返回品牌和 locale 信息
- 不携带 tenant 内部配置、provider 参数或 feature flag 明细
- 保留独立限流

### `/public/kb/*`

- 只暴露已公开发布的知识内容
- 结果集必须支持分页 / 限制条数
- 关注抓取、枚举和全文爬取风险
- 运行基线要求为 `/public/kb/` 配置独立路径级限流

### `/public/csat/*`

- 避免通过可预测 ID 枚举他人满意度数据
- 对提交频率做独立限制
- 运行基线要求为 `/public/csat/` 配置独立路径级限流

### `/api/v1/auth/*` 匿名入口

- `/login`、`/register`、`/refresh` 必须有独立路径级限流
- 不在错误响应里泄漏账号存在性、内部状态或策略细节
- 对暴力尝试、refresh 滥用、批量注册建立告警

### `/api/v1/ws`

- 匿名建连入口要控制建连频率和并发数
- 不应在握手阶段暴露内部状态、调度信息或租户敏感元数据
- 运行基线要求为 `/api/v1/ws` 配置独立路径级限流

### `/uploads/*`

- 只暴露确有公开需求的上传资产
- 不回显真实磁盘路径或内部存储 key 规则
- 运行基线要求为 `/uploads/` 配置独立路径级限流

## 与现有文档的关系

- 路由表面分类：[auth-surface-policy.md](auth-surface-policy.md)
- 运行安全基线：[security-baseline-operations.md](security-baseline-operations.md)
- token / key 轮换：[token-lifecycle-and-key-rotation.md](token-lifecycle-and-key-rotation.md)
