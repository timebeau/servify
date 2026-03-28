import { request } from '@umijs/max';

export async function getPortalConfig() {
  return request('/public/portal/config');
}
