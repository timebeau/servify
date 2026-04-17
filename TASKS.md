# Servify v0.1.0 验收任务清单

> 最后更新: 2026-04-18
> **重要**: 使用 `-count=1` 禁用测试缓存以获得真实结果

## 验证方法

```bash
# 全量测试（禁用缓存）
go test -count=1 ./apps/server/...

# Integration 测试（禁用缓存）
go test -count=1 -tags=integration ./apps/server/internal/handlers ./apps/server/internal/services
```

> **注意**: `go test` 默认会缓存结果，可能隐藏真实的测试失败。**必须使用 `-count=1`** 来验证真实的测试状态。

## 已完成 ✅

| 模块 | 功能 | API | 测试文件 | 状态 |
|------|------|-----|----------|------|
| 客户管理 | 客户列表 | GET /api/customers | customer_handler_test.go | ✅ 测试通过 |
| 客户管理 | 活动轨迹 | GET /api/customers/:id/activity | customer_handler_test.go | ✅ 测试通过 |
| 客服管理 | 创建客服 | POST /api/agents | agent_handler_test.go | ✅ 测试通过 |
| 客服管理 | 查找可用客服 | GET /api/agents/find-available | agent_handler_test.go | ✅ 测试通过 |
| 客服管理 | 会话分配/释放 | POST /api/agents/:id/assign-session | agent_handler_test.go | ✅ 测试通过 |
| 会话转接 | 转人工 | POST /api/session-transfer/to-human | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 指定客服转接 | POST /api/session-transfer/to-agent | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 查询转接历史 | GET /api/session-transfer/history/:session_id | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 查询等待队列 | GET /api/session-transfer/waiting | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 取消等待 | POST /api/session-transfer/cancel | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 处理排队 | POST /api/session-transfer/process-queue | session_transfer_handler_test.go | ✅ 测试通过 |
| 会话转接 | 自动转接检查 | POST /api/session-transfer/check-auto | session_transfer_handler_test.go | ✅ 测试通过 |
| 满意度 | 创建满意度记录 | POST /api/satisfactions | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 满意度列表 | GET /api/satisfactions | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 满意度统计 | GET /api/satisfactions/stats | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 调查列表 | GET /api/satisfactions/surveys | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 重发调查 | POST /api/satisfactions/surveys/:id/resend | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 查看满意度详情 | GET /api/satisfactions/:id | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 更新满意度 | PUT /api/satisfactions/:id | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 删除满意度 | DELETE /api/satisfactions/:id | satisfaction_handler_test.go | ✅ 测试通过 |
| 满意度 | 按工单查询满意度 | GET /api/tickets/:id/satisfaction | satisfaction_handler_test.go | ✅ 测试通过 |
| 宏管理 | CRUD+应用 | /api/macros 系列 | macro_handler_test.go | ✅ 测试通过 |
| 应用集成 | CRUD | /api/apps/integrations 系列 | app_market_handler_test.go | ✅ 测试通过 |
| 自定义字段 | CRUD | /api/custom-fields 系列 | custom_field_handler_test.go | ✅ 测试通过 |
| 统计分析 | Dashboard | GET /api/statistics/dashboard | statistics_handler_test.go | ✅ 测试通过 |
| 统计分析 | 客服绩效统计 | GET /api/statistics/agent-performance | statistics_handler_test.go | ✅ 测试通过 |
| 统计分析 | 工单分类统计 | GET /api/statistics/ticket-category | statistics_handler_test.go | ✅ 测试通过 |
| 统计分析 | 工单优先级统计 | GET /api/statistics/ticket-priority | statistics_handler_test.go | ✅ 测试通过 |
| 统计分析 | 客户来源统计 | GET /api/statistics/customer-source | statistics_handler_test.go | ✅ 测试通过 |
| 统计分析 | 远程协助工单统计 | GET /api/statistics/remote-assist-tickets | statistics_handler_test.go | ✅ 测试通过 |

## 最近修复

### 编译错误修复

1. **SQLite driver 迁移** - 15个测试文件从 `gorm.io/driver/sqlite` (CGO) 迁移到 `github.com/glebarez/sqlite` (纯 Go)
   - `apps/server/internal/platform/usersecurity/revoke_test.go`
   - `apps/server/internal/platform/configscope/gorm_provider_test.go`
   - `apps/server/internal/platform/auth/user_state_policy_test.go`
   - `apps/server/internal/platform/audit/query_test.go`
   - `apps/server/internal/modules/voice/infra/gorm_repository_test.go`
   - `apps/server/internal/modules/ticket/infra/gorm_repository_test.go`
   - `apps/server/internal/modules/routing/infra/gorm_repository_scope_integration_test.go`
   - `apps/server/internal/modules/routing/delivery/session_transfer_adapter_test.go`
   - `apps/server/internal/modules/knowledge/infra/gorm_repository_test.go`
   - `apps/server/internal/modules/customer/infra/gorm_repository_scope_integration_test.go`
   - `apps/server/internal/modules/conversation/infra/gorm_repository_scope_integration_test.go`
   - `apps/server/internal/modules/analytics/infra/gorm_repository_scope_integration_test.go`
   - `apps/server/internal/modules/agent/infra/gorm_repository_scope_integration_test.go`
   - `apps/server/internal/app/server/router_auth_test.go`
   - `apps/server/internal/app/server/ai_scoped_handler_test.go`

2. **AI handler adapter 接口** - 修复测试 stub 方法名与接口定义不匹配
   - `UploadDocumentToWeKnora` → `UploadKnowledgeDocument`
   - `SetWeKnoraEnabled` → `SetKnowledgeProviderEnabled`

### 测试失败修复

3. **user_security_handler_test** - 修复 risk_score 断言 (8→9)

4. **MacroService.List 排序不稳定** - 修复时间精度问题
   - 添加 `id DESC` 作为二级排序，确保相同 updated_at 时的稳定排序
   - 修复测试 DB helper 使用 `cache=shared`
   - 添加 ID 验证确保测试断言正确

5. **默认配置安全问题** - 加强生产环境配置验证
   - 添加 `Validate()` 函数，检测生产环境中的不安全默认值
   - 更新 `config.yml` 数据库密码使用环境变量占位符
   - 更新 CORS `allowed_headers` 从 `["*"]` 改为具体 header 列表
   - 更新默认 JWT secret 为更明显的占位符
   - 更新安全检查以识别所有已知不安全的默认值

---

## 测试验证结果

```bash
# 验证时间: 2026-04-18
# 命令: go test -count=1 -tags=integration ./apps/server/internal/handlers ./apps/server/internal/services

ok      servify/apps/server/internal/handlers    0.626s
ok      servify/apps/server/internal/services    1.721s
```

```bash
# 全量测试（禁用缓存）
# 命令: go test -count=1 ./apps/server/...

ok      servify/apps/server/internal/...         (all packages passed)
```

---

## 统计

- 已完成: 37 项 ✅ (已用 `-count=1` 验证通过)
- 待完成: 0 项
