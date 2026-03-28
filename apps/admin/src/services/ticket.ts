import { request } from '@umijs/max';

const API = '/api/tickets';

export async function listTickets(params: API.TicketListParams) {
  return request<API.PaginatedResponse<API.Ticket>>(API, { params });
}

export async function getTicket(id: number) {
  return request<API.Ticket>(`${API}/${id}`);
}

export async function createTicket(data: Partial<API.Ticket>) {
  return request<API.Ticket>(API, { method: 'POST', data });
}

export async function updateTicket(id: number, data: Partial<API.Ticket>) {
  return request<API.Ticket>(`${API}/${id}`, { method: 'PUT', data });
}

export async function assignTicket(id: number, agentId: number) {
  return request(`${API}/${id}/assign`, { method: 'POST', data: { agent_id: agentId } });
}

export async function closeTicket(id: number) {
  return request(`${API}/${id}/close`, { method: 'POST' });
}

export async function addComment(id: number, data: { content: string; internal?: boolean }) {
  return request(`${API}/${id}/comments`, { method: 'POST', data });
}

export async function getComments(id: number) {
  return request(`${API}/${id}/comments`);
}

export async function getTicketSatisfaction(id: number) {
  return request(`${API}/${id}/satisfaction`);
}

export async function bulkUpdateTickets(data: { ticket_ids: number[]; action: string }) {
  return request(`${API}/bulk`, { method: 'POST', data });
}

export async function exportTickets(params: API.TicketListParams) {
  return request(`${API}/export`, { params, responseType: 'blob' });
}
