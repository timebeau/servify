import { request } from '@umijs/max';

const API = '/api/automations';

export async function listAutomations(params?: { page?: number; page_size?: number }) {
  return request<API.PaginatedResponse<API.Automation>>(API, { params });
}

export async function createAutomation(data: Partial<API.Automation>) {
  return request<API.Automation>(API, { method: 'POST', data });
}

export async function deleteAutomation(id: number) {
  return request(`${API}/${id}`, { method: 'DELETE' });
}

export async function getAutomationRuns(params?: { page?: number; page_size?: number }) {
  return request(`${API}/runs`, { params });
}

export async function runAutomation(id: number) {
  return request(`${API}/run`, { method: 'POST', data: { automation_id: id } });
}

export async function bulkRunAutomations(data: { trigger_type: string; payload?: any }) {
  return request(`${API}/run`, { method: 'POST', data });
}
