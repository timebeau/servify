import { request } from '@/lib/request';

export async function getConversationMessages(sessionId: string) {
  return request<{ data: API.ConversationMessage[] }>(`/api/omni/sessions/${sessionId}/messages`);
}

export async function sendConversationMessage(sessionId: string, data: { content: string }) {
  return request(`/api/omni/sessions/${sessionId}/messages`, {
    method: 'POST',
    data,
  });
}
