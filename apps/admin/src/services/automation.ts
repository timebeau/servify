import { request } from '@/lib/request';

const API = '/api/automations';

export async function listAutomations(params?: { page?: number; page_size?: number }) {
  return request<API.PaginatedResponse<API.Automation>>(API, { params });
}

export async function createAutomation(data: Partial<API.Automation>) {
  return request<API.Automation>(API, { method: 'POST', data });
}

export async function deleteAutomation(id: number) {
  return request<API.MessageResponse>(`${API}/${id}`, { method: 'DELETE' });
}

export async function getAutomationRuns(params?: { page?: number; page_size?: number }) {
  return request<API.PaginatedResponse<API.AutomationRun>>(`${API}/runs`, { params });
}

export async function runAutomation(id: number) {
  // RunBatch accepts BatchRunRequest with event, ticket_ids, dry_run
  // For running a single automation by ID, use the event-based trigger
  return request<API.MessageResponse>(`${API}/run`, {
    method: 'POST',
    data: {
      event: 'manual_trigger',
      ticket_ids: [],
      dry_run: false,
      metadata: { automation_id: id },
    },
  });
}

export async function bulkRunAutomations(data: { trigger_type: string; payload?: Record<string, unknown> }) {
  return request<API.MessageResponse>(`${API}/run`, { method: 'POST', data });
}
