import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

const API = '/api/satisfactions';

export async function listSatisfactions(params?: { page?: number; page_size?: number; ticket_id?: number }) {
  const payload = await request<unknown>(API, { params });
  return normalizePaginatedResponse<API.Satisfaction>(payload);
}

export async function getSatisfactionStats() {
  return request<API.SatisfactionStats>(`${API}/stats`);
}

export async function listSurveys(params?: { page?: number; page_size?: number }) {
  const payload = await request<unknown>(`${API}/surveys`, { params });
  return normalizePaginatedResponse<API.SatisfactionSurvey>(payload);
}

export async function resendSurvey(id: number) {
  return request<API.MessageResponse>(`${API}/surveys/${id}/resend`, { method: 'POST' });
}
