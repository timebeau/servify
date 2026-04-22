import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

const API = '/api/custom-fields';

export async function listCustomFields(params?: {
  page?: number;
  page_size?: number;
  resource?: string;
  active?: boolean;
}) {
  const payload = await request<unknown>(API, {
    params: {
      resource: params?.resource,
      active: params?.active,
    },
  });
  return normalizePaginatedResponse<API.CustomField>(payload);
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
  return request<API.MessageResponse>(`${API}/${id}`, { method: 'DELETE' });
}
