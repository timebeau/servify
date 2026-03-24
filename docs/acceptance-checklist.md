# Servify 验收清单

这份清单的目标不是证明“看起来差不多做完了”，而是防止把“代码里有入口”误判成“功能已经可交付”。

验收时，每个功能项都必须同时具备以下 4 类证据，否则一律不算完成：

1. 代码入口：路由、命令或调用入口明确存在。
2. 自动化测试：至少有相关单测、集成测试或命令级验证。
3. 运行证据：服务真实启动，接口可访问，请求与响应符合预期。
4. 数据证据：涉及数据库、状态流转、文件或外部依赖时，能看到前后状态变化。

## 状态定义

- `未验`：只确认代码入口存在，没有实际跑通。
- `部分通过`：编译或单测通过，但关键业务链路、异常路径或数据落库未验。
- `通过`：具备代码、测试、运行、数据四类证据。
- `阻塞`：缺配置、依赖、权限、环境或数据，暂时无法验。

## 使用方式

先准备运行环境：

```bash
make migrate DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=password DB_NAME=servify
make run-cli CONFIG=./config.yml
make run-weknora CONFIG=./config.weknora.yml
```

常用基础验证命令：

```bash
make build
make build-weknora
make test
go test -tags weknora ./apps/server/cmd/cli ./apps/server/internal/handlers ./apps/server/pkg/weknora/...
```

每个条目都按下面模板记录：

| 字段 | 内容 |
| --- | --- |
| 功能项 | 例如“创建工单” |
| 入口 | 路由 / CLI / 页面 |
| 前置条件 | 数据库、配置、鉴权、依赖服务 |
| 操作步骤 | 明确到请求参数或操作动作 |
| 预期结果 | 状态码、返回体、状态变化 |
| 自动化证据 | 相关测试文件或命令 |
| 人工证据 | curl 结果、截图、日志、DB 查询结果 |
| 状态 | 未验 / 部分通过 / 通过 / 阻塞 |

## 当前已确认的最小事实

这些是当前已经实际验证过的，不代表“全量功能完成”：

| 项目 | 证据 | 状态 |
| --- | --- | --- |
| WeKnora CLI 可构建 | `make build-weknora` 通过 | 通过 |
| WeKnora 相关 handler/pkg 基础测试 | `go test -tags weknora ./apps/server/cmd/cli ./apps/server/internal/handlers ./apps/server/pkg/weknora/...` 通过 | 部分通过 |

说明：

- 上述结果只能证明构建链和部分测试是通的。
- 不能据此推出所有 API、权限、迁移、前端、数据链路都已完成。

## 一票否决项

以下任何一项失败，都不能对外宣称“已完成”：

- 服务无法启动。
- 数据库迁移失败。
- 核心路由返回 500 或返回结构与契约不符。
- 创建后无法查询，或更新后状态不变化。
- 鉴权与权限中间件失效。
- WeKnora 不可用时没有回退或错误不可控。
- 关键流程只有 happy path，没有异常路径验证。

## 全量验收矩阵

以下清单按代码实际暴露的功能面整理，主要依据：

- [README.md](/Users/cui/Workspaces/servify/README.md)
- [router.go](/Users/cui/Workspaces/servify/apps/server/internal/app/server/router.go)

### 1. 基础运行与可观测性

代码入口：

- [health.go](/Users/cui/Workspaces/servify/apps/server/internal/app/server/health.go)
- [run.go](/Users/cui/Workspaces/servify/apps/server/cmd/cli/run.go)
- [run_enhanced.go](/Users/cui/Workspaces/servify/apps/server/cmd/cli/run_enhanced.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 健康检查 | `GET /health` | 服务启动后直接请求 | `200` 且返回健康信息 | [health_enhanced_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/health_enhanced_test.go) | 未验 |
| 就绪检查 | `GET /ready` | 在依赖就绪后请求 | `200`；依赖异常时不应误报健康 | [health_enhanced_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/health_enhanced_test.go) | 未验 |
| Prometheus 指标 | `GET <metricsPath>` | 启动后访问 metrics | 返回文本指标；包含 runtime/AI/DB 指标 | [ai_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler.go) | 未验 |
| CLI 标准构建 | `make build` | 执行构建 | 生成二进制，无编译错误 | `make build` | 未验 |
| CLI WeKnora 构建 | `make build-weknora` | 执行构建 | 生成二进制，无编译错误 | `make build-weknora` | 通过 |
| 本地最小环境检查 | `make local-check` | 执行环境校验脚本 | 所有关键依赖通过 | [Makefile](/Users/cui/Workspaces/servify/Makefile) | 未验 |

### 2. 实时能力

代码入口：

- [router.go](/Users/cui/Workspaces/servify/apps/server/internal/app/server/router.go)
- [handlers.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/handlers.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| WebSocket 建连 | `GET /api/v1/ws` | 带 `session_id` 建立连接 | 成功升级协议，可收发消息 | [websocket_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/websocket_test.go) | 未验 |
| WebSocket 统计 | `GET /api/v1/ws/stats` | 建连前后各请求一次 | client 数量正确变化 | [websocket_processing_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/websocket_processing_test.go) | 未验 |
| WebRTC 统计 | `GET /api/v1/webrtc/stats` | 建立会话后请求 | 返回连接统计 | [webrtc_stats_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/webrtc_stats_test.go) | 未验 |
| WebRTC 连接列表 | `GET /api/v1/webrtc/connections` | 建立多连接后请求 | 返回当前连接列表 | [webrtc_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/webrtc_test.go) | 未验 |
| 平台消息路由统计 | `GET /api/v1/messages/platforms` | 模拟不同平台消息后请求 | 平台维度统计正确 | [router_stats_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/router_stats_test.go) | 未验 |

### 3. AI 与 WeKnora

代码入口：

- [router.go](/Users/cui/Workspaces/servify/apps/server/internal/app/server/router.go)
- [ai_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler.go)
- [WEKNORA_INTEGRATION.md](/Users/cui/Workspaces/servify/docs/WEKNORA_INTEGRATION.md)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| AI 查询 | `POST /api/v1/ai/query` | 发起普通问题请求 | 返回回答，结构稳定 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| AI 状态 | `GET /api/v1/ai/status` | 请求状态接口 | 返回 AI 服务当前状态 | [ai_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_test.go) | 未验 |
| AI 指标 | `GET /api/v1/ai/metrics` | 调用多次查询后请求 | query count、latency 等可见 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| 知识上传 | `POST /api/v1/ai/knowledge/upload` | 上传合法文档 | 上传成功并返回任务或文档信息 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| 知识同步 | `POST /api/v1/ai/knowledge/sync` | 触发同步 | 返回同步结果，无 500 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| 启用 WeKnora | `PUT /api/v1/ai/weknora/enable` | 启用增强能力 | 状态切换成功 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| 禁用 WeKnora | `PUT /api/v1/ai/weknora/disable` | 禁用增强能力 | 状态切换成功 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| 熔断器重置 | `POST /api/v1/ai/circuit-breaker/reset` | 先制造失败再重置 | 熔断状态恢复 | [ai_handler_comprehensive_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ai_handler_comprehensive_test.go) | 未验 |
| WeKnora 不可用回退 | AI 查询主链路 | 关闭 WeKnora 或制造超时 | 系统进入 fallback，而不是整体不可用 | [orchestrated_ai_fallback_integration_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/orchestrated_ai_fallback_integration_test.go) | 未验 |

### 4. 客户管理

代码入口：[customer_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 创建客户 | `POST /api/customers` | 提交合法客户资料 | 创建成功，返回 ID | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 客户列表 | `GET /api/customers` | 创建若干客户后查询 | 支持列表和筛选 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 客户统计 | `GET /api/customers/stats` | 准备不同客户样本后查询 | 返回统计信息正确 | [customer_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/customer_service_test.go) | 未验 |
| 客户详情 | `GET /api/customers/:id` | 查询已存在客户 | 返回正确详情 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 更新客户 | `PUT /api/customers/:id` | 更新字段后重新查询 | 数据被持久化 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 活动轨迹 | `GET /api/customers/:id/activity` | 产生客户相关活动后查询 | 轨迹完整、排序正确 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 添加备注 | `POST /api/customers/:id/notes` | 添加备注后重查 | 备注被记录 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |
| 更新标签 | `PUT /api/customers/:id/tags` | 更新标签后重查 | 标签覆盖或合并符合设计 | [customer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/customer_handler_test.go) | 未验 |

### 5. 客服管理

代码入口：[agent_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 创建客服 | `POST /api/agents` | 提交合法资料 | 创建成功 | [agent_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler_test.go) | 未验 |
| 客服列表 | `GET /api/agents` | 创建后查询 | 返回客服列表 | [agent_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler_test.go) | 未验 |
| 在线客服列表 | `GET /api/agents/online` | 切换在线离线状态后查询 | 仅返回在线客服 | [agent_handler_extended_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler_extended_test.go) | 未验 |
| 客服统计 | `GET /api/agents/stats` | 准备不同负载样本后查询 | 统计值正确 | [agent_service_more_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/agent_service_more_test.go) | 未验 |
| 查找可用客服 | `GET /api/agents/find-available` | 准备多个客服状态 | 返回最合适客服 | [agent_service_assignment_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/agent_service_assignment_test.go) | 未验 |
| 更新状态 | `PUT /api/agents/:id/status` | 更新在线/忙碌状态 | 状态持久化 | [agent_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler_test.go) | 未验 |
| 上下线切换 | `POST /api/agents/:id/online` `POST /api/agents/:id/offline` | 调用接口并重查 | 状态切换正确 | [agent_handler_extended_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/agent_handler_extended_test.go) | 未验 |
| 会话分配/释放 | `POST /api/agents/:id/assign-session` `POST /api/agents/:id/release-session` | 分配后查看客服负载 | 负载变化正确 | [agent_service_assignment_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/agent_service_assignment_test.go) | 未验 |

### 6. 工单

代码入口：[ticket_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 创建工单 | `POST /api/tickets` | 提交合法工单 | 创建成功，返回 ID | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |
| 批量更新工单 | `POST /api/tickets/bulk` | 批量更新多个工单 | 返回批量结果，数据一致 | [ticket_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler.go) | 未验 |
| 工单列表 | `GET /api/tickets` | 创建后查询 | 支持筛选和分页 | [query_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/ticket/application/query_service_test.go) | 未验 |
| 导出工单 CSV | `GET /api/tickets/export` | 准备样本后导出 | 下载成功，字段完整 | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |
| 工单统计 | `GET /api/tickets/stats` | 准备不同状态工单后查询 | 统计值正确 | [query_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/ticket/application/query_service_test.go) | 未验 |
| 工单详情 | `GET /api/tickets/:id` | 查询已存在工单 | 返回正确详情 | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |
| 更新工单 | `PUT /api/tickets/:id` | 更新标题、优先级、状态 | 数据被持久化 | [command_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/ticket/application/command_service_test.go) | 未验 |
| 指派工单 | `POST /api/tickets/:id/assign` | 指派给客服 | 负责人变化正确 | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |
| 添加评论 | `POST /api/tickets/:id/comments` | 添加评论后重查 | 评论可读出 | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |
| 关闭工单 | `POST /api/tickets/:id/close` | 关闭后重查 | 状态变为关闭，关闭时间正确 | [ticket_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/ticket_handler_test.go) | 未验 |

### 7. 会话转接

代码入口：[session_transfer_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 转人工 | `POST /api/session-transfer/to-human` | 发起 AI 到人工转接 | 进入等待或转接成功 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 指定客服转接 | `POST /api/session-transfer/to-agent` | 指定客服转接 | 记录目标客服并更新状态 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 查询转接历史 | `GET /api/session-transfer/history/:session_id` | 转接后查询历史 | 历史完整可追溯 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 查询等待队列 | `GET /api/session-transfer/waiting` | 制造排队后查询 | 返回等待记录 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 取消等待 | `POST /api/session-transfer/cancel` | 取消等待记录 | 队列移除成功 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 处理排队 | `POST /api/session-transfer/process-queue` | 触发处理 | 有可用客服时完成派发 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |
| 自动转接检查 | `POST /api/session-transfer/check-auto` | 构造触发条件后执行 | 符合策略的会话被处理 | [session_transfer_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/session_transfer_handler_test.go) | 未验 |

### 8. 满意度与公开问卷

代码入口：

- [satisfaction_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler.go)
- [csat_public_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/csat_public_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 创建满意度记录 | `POST /api/satisfactions` | 提交评价 | 创建成功 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 满意度列表 | `GET /api/satisfactions` | 创建后查询 | 返回列表和筛选结果 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 满意度统计 | `GET /api/satisfactions/stats` | 准备多样本后查询 | 聚合结果正确 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 调查列表 | `GET /api/satisfactions/surveys` | 查询调查任务 | 返回调查清单 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 重发调查 | `POST /api/satisfactions/surveys/:id/resend` | 重发问卷 | 任务状态更新 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 查看满意度详情 | `GET /api/satisfactions/:id` | 查询单条记录 | 返回正确详情 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 更新满意度 | `PUT /api/satisfactions/:id` | 修改评分或备注 | 数据更新成功 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 删除满意度 | `DELETE /api/satisfactions/:id` | 删除后重查 | 记录不可再读 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 按工单查询满意度 | `GET /api/tickets/:id/satisfaction` | 工单已有评价时查询 | 返回对应满意度 | [satisfaction_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/satisfaction_handler_test.go) | 未验 |
| 公开问卷查看 | `GET /public/csat/:token` | 使用有效 token 访问 | 返回问卷内容 | [csat_public_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/csat_public_handler.go) | 未验 |
| 公开问卷提交 | `POST /public/csat/:token/respond` | 提交问卷响应 | 返回成功并持久化 | [csat_public_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/csat_public_handler.go) | 未验 |

### 9. 工作台、宏、集成、自定义字段

代码入口：

- [workspace_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/workspace_handler.go)
- [macro_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/macro_handler.go)
- [app_market_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/app_market_handler.go)
- [custom_field_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/custom_field_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 工作台概览 | `GET /api/omni/workspace` | 请求工作台概览 | 返回聚合视图 | [workspace_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/workspace_handler_test.go) | 未验 |
| 宏列表/创建/更新/删除/应用 | `/api/macros` 系列 | 依次执行 CRUD 与 apply | 规则可持久化且可实际应用 | [macro_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/macro_handler_test.go) | 未验 |
| 应用集成列表/创建/更新/删除 | `/api/apps/integrations` 系列 | 执行 CRUD | 集成配置生效并可回读 | [app_market_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/app_market_handler_test.go) | 未验 |
| 自定义字段列表/详情/创建/更新/删除 | `/api/custom-fields` 系列 | 执行 CRUD | 字段定义持久化并可用于工单等实体 | [custom_field_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/custom_field_handler_test.go) | 未验 |

### 10. 统计、SLA、排班

代码入口：

- [statistics_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/statistics_handler.go)
- [sla_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler.go)
- [shift_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/shift_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 仪表盘统计 | `GET /api/statistics/dashboard` | 准备样本数据后查询 | 返回聚合概览 | [statistics_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/statistics_handler_test.go) | 未验 |
| 时间范围统计 | `GET /api/statistics/time-range` | 指定日期范围查询 | 返回按时间切片聚合 | [statistics_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/statistics_service_test.go) | 未验 |
| 客服绩效统计 | `GET /api/statistics/agent-performance` | 准备客服样本后查询 | 返回绩效排序 | [statistics_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/statistics_service_test.go) | 未验 |
| 工单分类统计 | `GET /api/statistics/ticket-category` | 准备不同分类工单 | 分类聚合正确 | [statistics_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/statistics_handler_test.go) | 未验 |
| 工单优先级统计 | `GET /api/statistics/ticket-priority` | 准备不同优先级工单 | 优先级聚合正确 | [statistics_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/statistics_handler_test.go) | 未验 |
| 客户来源统计 | `GET /api/statistics/customer-source` | 准备不同来源客户 | 来源聚合正确 | [statistics_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/statistics_handler_test.go) | 未验 |
| 每日统计更新 | `POST /api/statistics/update-daily` | 指定日期触发更新 | 更新成功且可被后续查询读到 | [statistics_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/statistics_service_test.go) | 未验 |
| SLA 配置 CRUD | `/api/sla/configs` 系列 | 执行配置增删改查 | 配置生效并可回读 | [sla_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_test.go) | 未验 |
| SLA 优先级配置查询 | `GET /api/sla/configs/priority/:priority` | 查询指定优先级 | 返回正确规则 | [sla_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_test.go) | 未验 |
| SLA 违约列表 | `GET /api/sla/violations` | 准备违约工单后查询 | 返回违约记录 | [sla_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_test.go) | 未验 |
| SLA 违约解决 | `POST /api/sla/violations/:id/resolve` | 解决违约后重查 | 状态变化正确 | [sla_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_test.go) | 未验 |
| SLA 统计 | `GET /api/sla/stats` | 准备样本后查询 | 返回 SLA 统计视图 | [sla_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_test.go) | 未验 |
| 工单 SLA 检查 | `POST /api/sla/check/ticket/:ticket_id` | 对指定工单执行检查 | 返回 SLA 判定结果 | [sla_handler_check_ticket_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/sla_handler_check_ticket_test.go) | 未验 |
| 排班 CRUD 与统计 | `/api/shifts` 系列 | 创建、查询、更新、删除、统计 | 排班记录和统计正确 | [shift_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/shift_handler_test.go) | 未验 |

### 11. 自动化、知识库、辅助建议、激励

代码入口：

- [automation_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/automation_handler.go)
- [knowledge_doc_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/knowledge_doc_handler.go)
- [suggestion_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/suggestion_handler.go)
- [gamification_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/gamification_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 自动化触发器列表/创建/删除 | `/api/automations` 系列 | 执行 CRUD | 触发器生效并可管理 | [automation_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/automation_handler_test.go) | 未验 |
| 自动化运行记录 | `GET /api/automations/runs` | 触发自动化后查询 | 可查看运行结果 | [automation_service_more_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/automation_service_more_test.go) | 未验 |
| 自动化批量运行 | `POST /api/automations/run` | 触发批量运行 | 返回执行结果，失败项可见 | [automation_service_more_test.go](/Users/cui/Workspaces/servify/apps/server/internal/services/automation_service_more_test.go) | 未验 |
| 知识文档 CRUD | `/api/knowledge-docs` 系列 | 执行 CRUD | 文档可持久化 | [knowledge_doc_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/knowledge_doc_handler_test.go) | 未验 |
| 公开知识库查询 | `/public/kb/docs` 系列 | 访问公开知识文档 | 可读且权限边界正确 | [knowledge_doc_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/knowledge_doc_handler_test.go) | 未验 |
| 辅助建议 | `GET /api/assist/suggest` `POST /api/assist/suggest` | 传入上下文获取建议 | 返回建议结果 | [suggestion_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/suggestion_handler_test.go) | 未验 |
| 激励排行 | `GET /api/gamification/leaderboard` | 准备积分样本后查询 | 排名正确 | [gamification_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/gamification_handler_test.go) | 未验 |

### 12. 语音

代码入口：[voice_handler.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/voice_handler.go)

| 功能项 | 入口 | 验收步骤 | 预期结果 | 自动化证据 | 状态 |
| --- | --- | --- | --- | --- | --- |
| 协议列表 | `GET /api/voice/protocols` | 请求协议列表 | 返回已注册协议 | [voice_handler_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/voice_handler_test.go) | 未验 |
| 协议信令事件 | `POST /api/voice/protocols/:protocol/call-events/:event` | 发送 invite/answer/hangup 等事件 | 事件被正确映射和处理 | [voice_handler_integration_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/voice_handler_integration_test.go) | 未验 |
| 协议媒体事件 | `POST /api/voice/protocols/:protocol/media-events/:event` | 发送媒体事件 | 媒体状态正确更新 | [voice_handler_integration_test.go](/Users/cui/Workspaces/servify/apps/server/internal/handlers/voice_handler_integration_test.go) | 未验 |
| 开始录音 | `POST /api/voice/recordings/start` | 启动录音 | 返回 recording ID | [service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/voice/application/service_test.go) | 未验 |
| 停止录音 | `POST /api/voice/recordings/stop` | 停止录音后查询 | 录音状态完成 | [recording_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/voice/application/recording_service_test.go) | 未验 |
| 获取录音 | `GET /api/voice/recordings/:recordingID` | 查询录音详情 | 返回录音元数据 | [recording_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/voice/application/recording_service_test.go) | 未验 |
| 转写追加 | `POST /api/voice/transcripts` | 追加语音转写片段 | 转写被保存 | [transcript_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/voice/application/transcript_service_test.go) | 未验 |
| 转写列表 | `GET /api/voice/transcripts` | 追加后查询 | 返回转写内容 | [transcript_service_test.go](/Users/cui/Workspaces/servify/apps/server/internal/modules/voice/application/transcript_service_test.go) | 未验 |

## 权限与异常路径必须额外验

即使主流程通过，也不能直接判定完成，还必须补下面这些：

- 未登录访问 `/api/*` 是否被正确拦截。
- 权限不足访问不同资源是否返回正确错误。
- 非法参数是否返回 `400` 而不是 `500`。
- 查不存在 ID 是否返回 `404` 或明确错误。
- 外部依赖失败时是否有降级、重试或熔断行为。
- 重复提交、并发操作、幂等等是否符合预期。

相关入口见 [router.go](/Users/cui/Workspaces/servify/apps/server/internal/app/server/router.go) 中的鉴权与资源权限中间件注册。

## 建议的验收执行顺序

1. 环境级：`make migrate`、`make build`、`make build-weknora`、服务启动。
2. 存活级：`/health`、`/ready`、metrics。
3. 核心业务级：客户、客服、工单、会话转接。
4. 增值能力级：AI、知识库、自动化、统计、SLA。
5. 扩展能力级：语音、公开问卷、公开知识库。
6. 负向与权限级：未授权、非法参数、外部依赖失效。

## 对外宣称“已完成”的最低门槛

至少满足以下条件，才能说“该模块完成验收”：

- 核心路由的 happy path 与异常路径都跑通。
- 自动化测试能覆盖该模块主要读写路径。
- 数据落库、状态流转、导出或副作用经过实际核验。
- 对应条目在本清单中被逐项标记为 `通过`，并附证据。

否则，最多只能说：

- “代码入口已存在”
- “构建已通过”
- “部分测试已通过”
- “功能仍待联调/待验收”

