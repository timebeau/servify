import { request } from '@/lib/request';

export async function getWorkspaceOverview() {
  return request('/api/omni/workspace');
}
