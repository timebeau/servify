import { request } from '@umijs/max';

const API = '/api/security';

export async function getUserSecurity(userId: number) {
  return request(`${API}/users/${userId}`);
}

export async function revokeUserTokens(userId: number) {
  return request(`${API}/users/${userId}/revoke-tokens`, { method: 'POST' });
}
