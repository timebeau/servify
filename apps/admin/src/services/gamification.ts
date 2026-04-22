import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

type LeaderboardEntry = Partial<API.LeaderboardRecord> & {
  agent_id?: number | string;
  agent_name?: string;
  username?: string;
  csat_avg?: number | string;
};

function toFiniteNumber(value: unknown, fallback = 0): number {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value);
    if (Number.isFinite(parsed)) {
      return parsed;
    }
  }
  return fallback;
}

function toDisplayName(entry: LeaderboardEntry): string {
  if (typeof entry.agent === 'string' && entry.agent.trim() !== '') {
    return entry.agent;
  }
  if (typeof entry.agent_name === 'string' && entry.agent_name.trim() !== '') {
    return entry.agent_name;
  }
  if (typeof entry.username === 'string' && entry.username.trim() !== '') {
    return entry.username;
  }
  return '-';
}

function normalizeLeaderboardRecord(entry: LeaderboardEntry, index: number): API.LeaderboardRecord {
  const agent = toDisplayName(entry);
  const id = entry.id ?? entry.agent_id ?? `${agent}-${index + 1}`;

  return {
    id,
    rank: toFiniteNumber(entry.rank, index + 1),
    agent,
    agent_id: entry.agent_id,
    department: typeof entry.department === 'string' ? entry.department : undefined,
    score: toFiniteNumber(entry.score),
    resolved_tickets: toFiniteNumber(entry.resolved_tickets),
    avg_rating: toFiniteNumber(entry.avg_rating ?? entry.csat_avg),
    avg_response_time: toFiniteNumber(entry.avg_response_time),
    badges: Array.isArray(entry.badges) ? entry.badges : undefined,
  };
}

export async function getLeaderboard(params?: {
  limit?: number;
  days?: number;
  start_date?: string;
  end_date?: string;
  department?: string;
}) {
  const payload = await request<unknown>('/api/gamification/leaderboard', { params });
  const normalized = normalizePaginatedResponse<LeaderboardEntry>(payload, ['entries', 'data', 'items']);

  return {
    ...normalized,
    total: normalized.total || normalized.data.length,
    data: normalized.data.map(normalizeLeaderboardRecord),
  };
}
