/** 工单状态 */
export const TICKET_STATUS_MAP: Record<string, { text: string; color: string }> = {
  open: { text: '待处理', color: 'blue' },
  assigned: { text: '已分配', color: 'cyan' },
  in_progress: { text: '处理中', color: 'orange' },
  pending_customer: { text: '待客户回复', color: 'gold' },
  resolved: { text: '已解决', color: 'green' },
  closed: { text: '已关闭', color: 'default' },
};

/** 工单优先级 */
export const TICKET_PRIORITY_MAP: Record<string, { text: string; color: string }> = {
  low: { text: '低', color: 'default' },
  medium: { text: '中', color: 'blue' },
  high: { text: '高', color: 'orange' },
  urgent: { text: '紧急', color: 'red' },
};

/** 客服在线状态 */
export const AGENT_STATUS_MAP: Record<string, { text: string; color: string }> = {
  online: { text: '在线', color: 'green' },
  busy: { text: '忙碌', color: 'orange' },
  away: { text: '离开', color: 'default' },
  offline: { text: '离线', color: 'default' },
};

/** 客户来源 */
export const CUSTOMER_SOURCE_MAP: Record<string, string> = {
  web: '网页',
  wechat: '微信',
  telegram: 'Telegram',
  api: 'API',
  import: '导入',
};

/** SLA 违约状态 */
export const SLA_VIOLATION_STATUS_MAP: Record<string, { text: string; color: string }> = {
  pending: { text: '待处理', color: 'red' },
  acknowledged: { text: '已确认', color: 'orange' },
  resolved: { text: '已解决', color: 'green' },
};
