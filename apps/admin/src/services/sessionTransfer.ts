import { request } from '@/lib/request';

export async function getWaitingQueue() {
  return request<{ count: number; data: API.TransferQueueRecord[] }>('/api/session-transfer/waiting');
}

export async function getTransferHistory(sessionId: string) {
  return request<API.SessionTransferRecord[]>(`/api/session-transfer/history/${sessionId}`);
}

export async function transferToAgent(data: { session_id: string; agent_id: number; reason?: string }) {
  return request<API.SessionTransferResult>('/api/session-transfer/to-agent', {
    method: 'POST',
    data: {
      session_id: data.session_id,
      target_agent_id: data.agent_id,
      reason: data.reason,
    },
  });
}

export async function cancelTransfer(sessionId: string) {
  return request<API.MessageResponse<{ session_id: string }>>('/api/session-transfer/cancel', {
    method: 'POST',
    data: { session_id: sessionId },
  });
}

export async function processQueue() {
  return request<API.MessageResponse>('/api/session-transfer/process-queue', { method: 'POST' });
}
