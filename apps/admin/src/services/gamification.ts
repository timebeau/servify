import { request } from '@/lib/request';

export async function getLeaderboard(params?: { page?: number; page_size?: number; period?: string }) {
  return request<API.LeaderboardRecord[] | { data: API.LeaderboardRecord[]; total?: number }>(
    '/api/gamification/leaderboard',
    { params },
  );
}
