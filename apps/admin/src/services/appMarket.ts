import { request } from '@/lib/request';

const API = '/api/integrations';

export async function listIntegrations(params?: { page?: number; page_size?: number; type?: string }) {
  return request<API.PaginatedResponse<API.Integration>>(API, { params });
}

export async function createIntegration(data: Partial<API.Integration>) {
  return request<API.Integration>(API, { method: 'POST', data });
}

export async function updateIntegration(id: number, data: Partial<API.Integration>) {
  return request<API.Integration>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteIntegration(id: number) {
  return request<API.MessageResponse>(`${API}/${id}`, { method: 'DELETE' });
}
