import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

const API = '/api/sla';

export async function listSLAConfigs(params?: { page?: number; page_size?: number }) {
  const payload = await request<unknown>(`${API}/configs`, { params });
  return normalizePaginatedResponse<API.SLAConfig>(payload);
}

export async function createSLAConfig(data: Partial<API.SLAConfig>) {
  return request<API.SLAConfig>(`${API}/configs`, { method: 'POST', data });
}

export async function updateSLAConfig(id: number, data: Partial<API.SLAConfig>) {
  return request<API.SLAConfig>(`${API}/configs/${id}`, { method: 'PUT', data });
}

export async function deleteSLAConfig(id: number) {
  return request<API.MessageResponse>(`${API}/configs/${id}`, { method: 'DELETE' });
}

export async function getSLAStats() {
  return request<API.SLAStats>(`${API}/stats`);
}

export async function listSLAViolations(params?: { page?: number; page_size?: number; status?: string }) {
  const payload = await request<unknown>(`${API}/violations`, {
    params: {
      page: params?.page,
      page_size: params?.page_size,
      resolved:
        params?.status === 'resolved'
          ? true
          : params?.status === 'pending'
            ? false
            : undefined,
    },
  });
  return normalizePaginatedResponse<API.SLAViolation>(payload);
}

export async function resolveSLAViolation(id: number) {
  return request<API.MessageResponse>(`${API}/violations/${id}/resolve`, { method: 'POST' });
}
