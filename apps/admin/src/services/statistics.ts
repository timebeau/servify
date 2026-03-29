import { request } from '@/lib/request';

/** 仪表板统计 */
export async function getDashboardStats() {
  return request<API.DashboardStats>('/api/statistics/dashboard');
}

/** 时间范围统计 */
export async function getTimeRangeStats(params: {
  start: string;
  end: string;
  interval?: string;
}) {
  return request('/api/statistics/time-range', { params });
}

/** 工单分类统计 */
export async function getTicketCategoryStats() {
  return request('/api/statistics/ticket-category');
}

/** 工单优先级统计 */
export async function getTicketPriorityStats() {
  return request('/api/statistics/ticket-priority');
}

/** 客服绩效 */
export async function getAgentPerformanceStats() {
  return request('/api/statistics/agent-performance');
}

/** 客户来源分布 */
export async function getCustomerSourceStats() {
  return request('/api/statistics/customer-source');
}
