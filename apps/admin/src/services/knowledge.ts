import { request } from '@/lib/request';

const API = '/api/knowledge-docs';

export async function listDocs(params: {
  page?: number;
  page_size?: number;
  search?: string;
  category?: string;
  status?: string;
}) {
  return request<API.PaginatedResponse<API.KnowledgeDoc>>(API, { params });
}

export async function getDoc(id: number) {
  return request<API.KnowledgeDoc>(`${API}/${id}`);
}

export async function createDoc(data: Partial<API.KnowledgeDoc>) {
  return request<API.KnowledgeDoc>(API, { method: 'POST', data });
}

export async function updateDoc(id: number, data: Partial<API.KnowledgeDoc>) {
  return request<API.KnowledgeDoc>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteDoc(id: number) {
  return request(`${API}/${id}`, { method: 'DELETE' });
}
