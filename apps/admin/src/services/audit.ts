import { request } from '@/lib/request';

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
  from?: string;
  to?: string;
}) {
  return request<API.PaginatedResponse<API.AuditLog>>(`${API}/logs`, {
    params: {
      page: params?.page,
      page_size: params?.page_size,
      action: params?.action,
      resource_type: params?.resource_type,
      resource_id: params?.resource_id,
      principal_kind: params?.principal_kind,
      actor_user_id: params?.actor_user_id,
      success: params?.success,
      from: params?.from,
      to: params?.to,
    },
  });
}
