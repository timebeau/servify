import { request } from '@umijs/max';

const API = '/api/voice';

export async function listProtocols(params?: { page?: number; page_size?: number }) {
  return request<API.PaginatedResponse<API.VoiceProtocol>>(`${API}/protocols`, { params });
}

export async function startRecording(protocolId: string) {
  return request(`${API}/recordings/start`, { method: 'POST', data: { protocol_id: protocolId } });
}

export async function stopRecording(recordingId: string) {
  return request(`${API}/recordings/stop`, { method: 'POST', data: { recording_id: recordingId } });
}

export async function getRecording(recordingId: string) {
  return request(`${API}/recordings/${recordingId}`);
}

export async function listTranscripts(params?: { page?: number; page_size?: number }) {
  return request(`${API}/transcripts`, { params });
}
