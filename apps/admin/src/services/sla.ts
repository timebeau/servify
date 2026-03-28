import { request } from '@umijs/max';

const API = '/api/sla';

export async function listSLAConfigs(params?: { page?: number; page_size?: number }) {
  return request<API.PaginatedResponse<API.SLAConfig>>(`${API}/configs`, { params });
}

export async function createSLAConfig(data: Partial<API.SLAConfig>) {
  return request(`${API}/configs`, { method: 'POST', data });
}

export async function updateSLAConfig(id: number, data: Partial<API.SLAConfig>) {
  return request(`${API}/configs/${id}`, { method: 'PUT', data });
}

export async function deleteSLAConfig(id: number) {
  return request(`${API}/configs/${id}`, { method: 'DELETE' });
}

export async function getSLAStats() {
  return request(`${API}/stats`);
}

export async function listSLAViolations(params?: { page?: number; page_size?: number; status?: string }) {
  return request<API.PaginatedResponse<any>>(`${API}/violations`, { params });
}

export async function resolveSLAViolation(id: number) {
  return request(`${API}/violations/${id}/resolve`, { method: 'POST' });
}
