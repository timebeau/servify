import { request } from '@/lib/request';

const API = '/api/custom-fields';

export async function listCustomFields(params?: { page?: number; page_size?: number; entity_type?: string }) {
  return request<API.PaginatedResponse<API.CustomField>>(API, { params });
}

export async function getCustomField(id: number) {
  return request<API.CustomField>(`${API}/${id}`);
}

export async function createCustomField(data: Partial<API.CustomField>) {
  return request<API.CustomField>(API, { method: 'POST', data });
}

export async function updateCustomField(id: number, data: Partial<API.CustomField>) {
  return request<API.CustomField>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteCustomField(id: number) {
  return request(`${API}/${id}`, { method: 'DELETE' });
}
