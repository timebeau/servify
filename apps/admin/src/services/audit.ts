import { request } from '@umijs/max';

const API = '/api/audit';

export async function listAuditLogs(params?: {
  page?: number;
  page_size?: number;
  action?: string;
  resource_type?: string;
  resource_id?: string;
  principal_kind?: string;
  actor_user_id?: number;
  success?: boolean;
  start_date?: string;
  end_date?: string;
}) {
  return request<API.PaginatedResponse<API.AuditLog>>(`${API}/logs`, { params });
}
