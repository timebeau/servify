# Auth Surface Policy

本文定义 Servify HTTP 路由的授权表面分类，目标是把“这个入口到底该不该匿名、该不该走 management token、该不该只允许 service”从零散 handler 逻辑里收口到统一规则。

## 分类

### 1. Public Surface

适用范围：

- 面向匿名访问者、浏览器或终端用户的公开读接口
- 只暴露最小只读能力或建连能力

约束：

- 不依赖 management JWT
- 不承载后台管理写操作
- 必须有独立的安全评审和路径级限流基线

当前路由：

- `/health`
- `/ready`
- `/public/portal/config`
- `/public/kb/*`
- `/public/csat/*`
- `/api/v1/ws`
- `/uploads/*`

### 2. Auth Surface

适用范围：

- 认证与会话自助相关入口
- 允许匿名进入认证前置流程，或允许已登录用户管理自己的 auth session

约束：

- 匿名入口必须具备独立路径级限流
- 已登录自助入口只允许操作自己的认证状态，不复用 management surface 的资源写权限
- 不承载跨用户、跨租户的后台管理动作

当前路由：

- 匿名认证入口：`/api/v1/auth/login`、`/api/v1/auth/register`、`/api/v1/auth/refresh`
- 已登录自助入口：`/api/v1/auth/me`、`/api/v1/auth/sessions`、`/api/v1/auth/sessions/logout-current`、`/api/v1/auth/sessions/logout-others`

### 3. Management Surface

适用范围：

- 管理后台、客服工作台、运营查询、人工触发的后台动作

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

### 4. Service Surface

适用范围：

- 机器到机器调用、采集上报、受控集成入口

约束：

- 必须经过 `AuthMiddleware`
- principal kind 仅允许 `service`
- 不复用 management `agent/admin` token
- 应优先设计成窄契约、低权限、可独立限流的入口

当前路由：

- `/api/v1/metrics/ingest`

## 新增路由准入规则

新增 HTTP 路由时，先回答这 4 个问题，再决定挂到哪个 surface：

1. 调用方是谁：匿名用户、终端用户、客服/管理员，还是服务账号？
2. 这个入口是公开读取、认证前置、后台管理，还是机器上报？
3. 失败时应该返回 `401`、`403`，还是允许匿名继续进入业务校验？
4. 是否需要专属路径级 rate limit，而不是只依赖全局限流？

默认策略：

- 拿不准时，不要直接挂到匿名 `public` surface
- 认证相关入口优先进入 `auth` surface，而不是混进 `public` 或 `management`
- 管理动作优先进入 `management` surface
- 机器上报、回调、采集优先进入 `service` surface

## 实现落点

- 路由装配：[apps/server/internal/app/server/router.go](/D:/workspaces/servify/apps/server/internal/app/server/router.go)
- 路由安全表面目录：[apps/server/internal/app/server/security_surface.go](/D:/workspaces/servify/apps/server/internal/app/server/security_surface.go)
- JWT 与 claims 归一化：[apps/server/internal/platform/auth](/D:/workspaces/servify/apps/server/internal/platform/auth)
- Gin 兼容入口：[apps/server/internal/middleware](/D:/workspaces/servify/apps/server/internal/middleware)

## 后续扩展

- 如果未来出现真正的 end-user authenticated business API，应新增独立 surface，而不是回退复用 management surface
- 如果引入回调签名、API key 或 mTLS，可在 service surface 内继续细分
- 如果开放接口继续增多，应把 `security surface catalog` 扩展成自动审查和审批流的输入源
