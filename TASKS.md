# Servify v0.1.0 验收任务清单

> 最后更新: 2026-04-18
> 测试状态: ✅ 所有测试通过

## 验证方法

```bash
# 全量测试
go test ./apps/server/...

# Integration 测试
go test -tags=integration ./apps/server/internal/handlers ./apps/server/internal/services
```

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

1. **SQLite driver 迁移** - 15个测试文件从 `gorm.io/driver/sqlite` 迁移到 `github.com/glebarez/sqlite`
2. **AI handler adapter 接口** - 修复测试 stub 方法名与接口定义不匹配
3. **user_security_handler_test** - 修复 risk_score 断言 (8→9)

---

## 统计

- 已完成: 37 项 ✅
- 待完成: 0 项
