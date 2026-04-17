# Servify v0.1.0 验收任务清单

> 生成时间: 2026-04-18
> 状态: 进行中

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

## 待完成任务 ⏳

### 1. 排班管理 (Shift)

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 排班 CRUD | /api/shifts 系列 | shift_handler_test.go | ⏳ 待验证 |
| 排班统计 | GET /api/shifts/stats | shift_handler_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] 创建排班记录成功
- [ ] 查询排班列表返回正确数据
- [ ] 更新排班记录生效
- [ ] 删除排班记录后不可再读
- [ ] 统计数据正确

---

### 2. 自动化 (Automation)

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 自动化触发器 CRUD | /api/automations 系列 | automation_handler_test.go | ⏳ 待验证 |
| 自动化运行记录 | GET /api/automations/runs | automation_service_more_test.go | ⏳ 待验证 |
| 自动化批量运行 | POST /api/automations/run | automation_service_more_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] 创建触发器成功
- [ ] 触发器正确执行
- [ ] 运行记录可查询
- [ ] 批量运行返回正确结果

---

### 3. 知识文档 (Knowledge Docs)

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 知识文档 CRUD | /api/knowledge-docs 系列 | knowledge_doc_handler_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] 创建文档成功
- [ ] 查询文档列表返回正确数据
- [ ] 更新文档生效
- [ ] 删除文档后不可再读

---

### 4. 辅助建议 (Assist Suggest)

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 辅助建议 | GET /api/assist/suggest | suggestion_handler_test.go | ⏳ 待验证 |
| 辅助建议 | POST /api/assist/suggest | suggestion_handler_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] GET 请求返回建议结果
- [ ] POST 请求返回建议结果
- [ ] 传入上下文后返回正确建议

---

### 5. 激励排行 (Gamification)

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 激励排行 | GET /api/gamification/leaderboard | gamification_handler_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] 准备积分样本后查询
- [ ] 排名正确显示
- [ ] 排序逻辑正确

---

### 6. 每日统计更新

| 功能 | API | 测试文件 | 状态 |
|------|-----|----------|------|
| 每日统计更新 | POST /api/statistics/update-daily | statistics_service_test.go | ⏳ 待验证 |

**验收标准:**
- [ ] 指定日期触发更新成功
- [ ] 更新后可被后续查询读到
- [ ] 聚合数据正确

---

### 7. 文档更新

| 文档 | 状态 |
|------|------|
| acceptance-checklist.md | ⏳ 待更新所有"未验"为"通过" |
| release-notes-v0.1.0.md | ⏳ 待最终确认 |

---

## 完成要求

每个任务需要:
1. **编写/运行测试用例** - 确保功能正常工作
2. **完成报告** - 记录测试结果和截图
3. **更新文档** - 更新 acceptance-checklist.md 状态

---

## 统计

- 已完成: 30 项
- 待完成: 7 大模块，约 13 项功能
