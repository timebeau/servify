import { request } from '@umijs/max';

export async function getWaitingQueue() {
  return request('/api/session-transfer/waiting');
}

export async function getTransferHistory(sessionId: string) {
  return request(`/api/session-transfer/history/${sessionId}`);
}

export async function transferToAgent(data: { session_id: string; agent_id: number; reason?: string }) {
  return request('/api/session-transfer/to-agent', { method: 'POST', data });
}

export async function cancelTransfer(sessionId: string) {
  return request('/api/session-transfer/cancel', { method: 'POST', data: { session_id: sessionId } });
}

export async function processQueue() {
  return request('/api/session-transfer/process-queue', { method: 'POST' });
}
