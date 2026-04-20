import { request } from '@/lib/request';

const API = '/api/agents';

export async function listAgents(params: {
  page?: number;
  page_size?: number;
  status?: string;
  search?: string;
}) {
  return request<API.PaginatedResponse<API.Agent>>(API, { params });
}

export async function getAgent(id: number) {
  return request<API.Agent>(`${API}/${id}`);
}

export async function createAgent(data: {
  name: string;
  email: string;
  skills?: string[];
  max_concurrent?: number;
}) {
  return request<API.Agent>(API, { method: 'POST', data });
}

export async function updateAgent(id: number, data: {
  name?: string;
  email?: string;
  skills?: string[];
  max_concurrent?: number;
}) {
  return request<API.Agent>(`${API}/${id}`, { method: 'PUT', data });
}

export async function updateAgentStatus(id: number, status: string) {
  return request<API.MessageResponse>(`${API}/${id}/status`, { method: 'PUT', data: { status } });
}

export async function agentOnline(id: number) {
  return request<API.MessageResponse>(`${API}/${id}/online`, { method: 'POST' });
}

export async function agentOffline(id: number) {
  return request<API.MessageResponse>(`${API}/${id}/offline`, { method: 'POST' });
}

export async function assignSessionToAgent(id: number, sessionId: string) {
  return request<API.MessageResponse>(`${API}/${id}/assign-session`, { method: 'POST', data: { session_id: sessionId } });
}

export async function releaseAgentSession(id: number, sessionId: string) {
  return request<API.MessageResponse>(`${API}/${id}/release-session`, { method: 'POST', data: { session_id: sessionId } });
}

export async function getOnlineAgents() {
  return request<API.ListResponse<API.Agent>>(`${API}/online`);
}

export async function getAgentStats() {
  return request<{
    total: number;
    online: number;
    busy: number;
    avg_response_time: number;
    avg_rating: number;
  }>(`${API}/stats`);
}
