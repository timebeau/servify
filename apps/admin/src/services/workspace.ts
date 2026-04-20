import { request } from '@/lib/request';

export async function getWorkspaceOverview() {
  return request<API.WorkspaceOverview>('/api/omni/workspace');
}
