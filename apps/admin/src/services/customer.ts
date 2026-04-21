import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

const API = '/api/customers';

export async function listCustomers(params: {
  page?: number;
  page_size?: number;
  search?: string;
  source?: string;
  industry?: string;
  tags?: string;
}) {
  const payload = await request<unknown>(API, { params });
  return normalizePaginatedResponse<API.Customer>(payload);
}

export async function getCustomer(id: number) {
  return request<API.Customer>(`${API}/${id}`);
}

export async function createCustomer(data: Partial<API.Customer>) {
  return request<API.Customer>(API, { method: 'POST', data });
}

export async function updateCustomer(id: number, data: Partial<API.Customer>) {
  return request<API.Customer>(`${API}/${id}`, { method: 'PUT', data });
}

export async function updateCustomerTags(id: number, tags: string[]) {
  return request<API.MessageResponse>(`${API}/${id}/tags`, { method: 'PUT', data: { tags } });
}

export async function addCustomerNote(id: number, content: string) {
  return request<API.MessageResponse>(`${API}/${id}/notes`, { method: 'POST', data: { content } });
}

export async function getCustomerActivity(id: number) {
  return request<API.CustomerActivity>(`${API}/${id}/activity`);
}

export async function revokeCustomerTokens(id: number) {
  return request<API.MessageResponse>(`${API}/${id}/revoke-tokens`, { method: 'POST' });
}

export async function getCustomerStats() {
  return request<{
    total_customers: number;
    active_customers: number;
    new_customers_today: number;
  }>(`${API}/stats`);
}
