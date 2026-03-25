# Auth Surface Policy

本文件定义 Servify HTTP 路由的授权表面分类，目标是避免把“是否需要 token”散落在各个 handler 里临时决定。

## 分类

### 1. Public Surface

适用范围：

- 面向终端用户、匿名访问者或浏览器入口
- 例如健康检查、公开问卷、公开知识库、公开 Portal 配置、匿名 realtime 建连入口

约束：

- 不依赖管理端 JWT
- 只允许暴露必要的只读或会话接入能力
- 不允许承载后台管理写操作

当前路由：

- `/health`
- `/ready`
- `/public/csat/*`
- `/public/kb/*`
- `/public/portal/config`
- `/api/v1/ws`

### 2. Management Surface

适用范围：

- 管理后台、客服工作台、运营查看、人工触发的后台动作
- 需要明确的 principal kind 与 RBAC 权限

约束：

- 必须经过 `AuthMiddleware`
- principal kind 仅允许 `agent`、`admin`、`service`
- 资源级访问继续叠加 `RequireResourcePermission(...)`
- 不允许 `end_user` token 进入该表面

当前路由：

- `/api/*`
- `/api/v1/ws/stats`
- `/api/v1/webrtc/*`
- `/api/v1/messages/platforms`
- `/api/v1/ai/*`

### 3. Service Surface

适用范围：

- 机器到机器调用、采集上报、内部 worker 或受控集成入口

约束：

- 必须经过 `AuthMiddleware`
- principal kind 仅允许 `service`
- 默认不复用管理端 `agent/admin` token
- 应优先设计成窄契约、低权限、可独立限流的入口

当前路由：

- `/api/v1/metrics/ingest`

## 新增路由准入规则

新增 HTTP 路由时，必须先回答下面 3 个问题，再决定挂到哪个 group：

1. 调用方是谁：匿名用户、终端用户、客服/管理员，还是服务账号？
2. 这个入口是业务管理动作、公共读取/接入，还是机器到机器上报？
3. 失败时应该返回 `401`、`403`，还是允许匿名继续进入业务校验？

默认策略：

- 拿不准时，不要直接挂到匿名 public surface
- 管理动作优先进入 management surface
- 机器上报、回调、采集优先进入 service surface

## 实现落点

- 路由装配：`apps/server/internal/app/server/router.go`
- JWT 与 claims 归一化：`apps/server/internal/platform/auth`
- Gin 兼容入口：`apps/server/internal/middleware`

## 后续扩展

- 若未来出现真正的 end-user authenticated API，应新增独立 surface，而不是复用 management surface
- 若引入内部回调签名、API key 或 mTLS，可在 service surface 上进一步细分，而不应回退到匿名入口
