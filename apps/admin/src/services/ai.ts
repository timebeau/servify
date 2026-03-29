import { request } from '@/lib/request';

export async function queryAI(data: { query: string; conversation_id?: string; customer_id?: number }) {
  return request('/api/v1/ai/query', { method: 'POST', data });
}

export async function getAIStatus() {
  return request<API.AIStatus>('/api/v1/ai/status');
}

export async function getAIMetrics() {
  return request('/api/v1/ai/metrics');
}
