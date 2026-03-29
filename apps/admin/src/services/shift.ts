import { request } from '@/lib/request';

const API = '/api/shifts';

export async function listShifts(params?: { page?: number; page_size?: number; agent_id?: number }) {
  return request<API.PaginatedResponse<API.Shift>>(API, { params });
}

export async function createShift(data: Partial<API.Shift>) {
  return request<API.Shift>(API, { method: 'POST', data });
}

export async function updateShift(id: number, data: Partial<API.Shift>) {
  return request<API.Shift>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteShift(id: number) {
  return request(`${API}/${id}`, { method: 'DELETE' });
}

export async function getShiftStats() {
  return request(`${API}/stats`);
}
