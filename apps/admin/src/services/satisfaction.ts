import { request } from '@/lib/request';

const API = '/api/satisfactions';

export async function listSatisfactions(params?: { page?: number; page_size?: number; ticket_id?: number }) {
  return request<API.PaginatedResponse<API.Satisfaction>>(API, { params });
}

export async function getSatisfactionStats() {
  return request(`${API}/stats`);
}

export async function listSurveys(params?: { page?: number; page_size?: number }) {
  return request(`${API}/surveys`, { params });
}

export async function resendSurvey(id: number) {
  return request(`${API}/surveys/${id}/resend`, { method: 'POST' });
}
