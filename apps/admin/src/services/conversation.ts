import { request } from '@/lib/request';

export async function getConversation(sessionId: string) {
  return request<{ data: API.Conversation }>(`/api/omni/sessions/${sessionId}`);
}

export async function getConversationMessages(
  sessionId: string,
  params?: { limit?: number; before?: string },
) {
  const query = new URLSearchParams();
  if (params?.limit) query.set('limit', String(params.limit));
  if (params?.before) query.set('before', params.before);
  const qs = query.toString();
  return request<{ data: API.ConversationMessage[] }>(
    `/api/omni/sessions/${sessionId}/messages${qs ? `?${qs}` : ''}`,
  );
}

export async function sendConversationMessage(sessionId: string, data: { content: string }) {
  return request(`/api/omni/sessions/${sessionId}/messages`, {
    method: 'POST',
    data,
  });
}

export async function assignAgent(sessionId: string, data: { agent_id: number }) {
  return request(`/api/omni/sessions/${sessionId}/assign`, {
    method: 'POST',
    data,
  });
}

export async function transferConversation(sessionId: string, data: { to_agent_id: number }) {
  return request(`/api/omni/sessions/${sessionId}/transfer`, {
    method: 'POST',
    data,
  });
}

export async function closeConversation(sessionId: string) {
  return request(`/api/omni/sessions/${sessionId}/close`, {
    method: 'POST',
  });
}
