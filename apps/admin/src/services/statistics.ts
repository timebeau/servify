import { request } from '@/lib/request';

/** 仪表板统计 */
export async function getDashboardStats() {
  return request<API.DashboardStats>('/api/statistics/dashboard');
}

/** 时间范围统计 */
export async function getTimeRangeStats(params: {
  start_date: string;
  end_date: string;
}) {
  return request('/api/statistics/time-range', { params });
}

/** 工单分类统计 */
export async function getTicketCategoryStats(params?: {
  start_date?: string;
  end_date?: string;
}) {
  return request('/api/statistics/ticket-category', { params });
}

/** 工单优先级统计 */
export async function getTicketPriorityStats(params?: {
  start_date?: string;
  end_date?: string;
}) {
  return request('/api/statistics/ticket-priority', { params });
}

/** 客服绩效 */
export async function getAgentPerformanceStats(params: {
  start_date: string;
  end_date: string;
  limit?: number;
}) {
  return request('/api/statistics/agent-performance', { params });
}

/** 客户来源分布 */
export async function getCustomerSourceStats() {
  return request('/api/statistics/customer-source');
}
