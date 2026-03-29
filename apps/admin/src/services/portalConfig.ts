import { request } from '@/lib/request';

export async function getPortalConfig() {
  return request('/public/portal/config');
}
