import { request } from '@umijs/max';

export async function getWorkspaceOverview() {
  return request('/api/omni/workspace');
}
