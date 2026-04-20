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

export async function createDoc(data: {
  title: string;
  content: string;
  category?: string;
  tags?: string[];
  is_public?: boolean;
}) {
  return request<API.KnowledgeDoc>(API, { method: 'POST', data });
}

export async function updateDoc(id: number, data: {
  title?: string;
  content?: string;
  category?: string;
  tags?: string[];
  is_public?: boolean;
}) {
  return request<API.KnowledgeDoc>(`${API}/${id}`, { method: 'PUT', data });
}

export async function deleteDoc(id: number) {
  return request<API.MessageResponse>(`${API}/${id}`, { method: 'DELETE' });
}
