import { request } from '@/lib/request';

export async function getLeaderboard(params?: { page?: number; page_size?: number; period?: string }) {
  return request('/api/gamification/leaderboard', { params });
}
