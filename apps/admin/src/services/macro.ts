import { request } from '@umijs/max';

const API = '/api/macros';

export async function listMacros(params?: { page?: number; page_size?: number; category?: string }) {
  return request<API.PaginatedResponse<API.Macro>>(API, { params });
}

export async function createMacro(data: Partial<API.Macro>) {
  return request<API.Macro>(API, { method: 'POST', data });
}

export async function updateMacro(id: number, data: Partial<API.Macro>) {
  return request<API.Macro>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteMacro(id: number) {
  return request(`${API}/${id}`, { method: 'DELETE' });
}

export async function applyMacro(id: number, data?: { ticket_id?: number; conversation_id?: string }) {
  return request(`${API}/${id}/apply`, { method: 'POST', data });
}
