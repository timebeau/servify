import { request } from '@/lib/request';

const API = '/api/security';

export async function getUserSecurity(userId: number) {
  return request<API.UserSecurityDetail>(`${API}/users/${userId}`);
}

export async function revokeUserTokens(userId: number) {
  return request(`${API}/users/${userId}/revoke-tokens`, { method: 'POST' });
}

export async function batchRevokeUserTokens(userIds: number[]) {
  return request<{ count: number; items: Array<{ user_id: number; token_version: number }> }>(
    `${API}/users/revoke-tokens`,
    {
      method: 'POST',
      data: { user_ids: userIds },
    },
  );
}

export async function queryUsersSecurity(userIds: number[]) {
  return request<{ count: number; items: API.UserSecurityPreview[] }>(`${API}/users/query`, {
    method: 'POST',
    data: { user_ids: userIds },
  });
}

export async function listUserSessions(userId: number) {
  return request<{ user_id: number; count: number; items: API.UserSecuritySession[] }>(
    `${API}/users/${userId}/sessions`,
  );
}

export async function revokeUserSession(userId: number, sessionId: string) {
  return request(`${API}/users/${userId}/sessions/revoke`, {
    method: 'POST',
    data: { session_id: sessionId },
  });
}

export async function revokeAllUserSessions(userId: number, exceptSessionId?: string) {
  return request<{ count: number; items: API.UserSecuritySession[] }>(
    `${API}/users/${userId}/sessions/revoke-all`,
    {
      method: 'POST',
      data: exceptSessionId ? { except_session_id: exceptSessionId } : {},
    },
  );
}

export async function revokeSecurityToken(token: string, reason?: string) {
  return request<API.RevokedToken & { message: string }>(`${API}/tokens/revoke`, {
    method: 'POST',
    data: { token, reason },
  });
}

export async function listRevokedTokens(params?: {
  page?: number;
  page_size?: number;
  jti?: string;
  user_id?: number | string;
  session_id?: string;
  token_use?: string;
  active_only?: boolean | string;
}) {
  return request<{ count: number; items: API.RevokedToken[] }>(`${API}/tokens/revoked`, {
    params,
  });
}
