import { request } from '@/lib/request';
import { normalizePaginatedResponse } from './_response';

const API = '/api/voice';

export async function listProtocols(params?: { page?: number; page_size?: number }) {
  const payload = await request<unknown>(`${API}/protocols`, { params });
  return normalizePaginatedResponse<API.VoiceProtocol>(payload);
}

export async function startRecording(callId: string) {
  return request<API.DataResponse<API.VoiceRecording>>(`${API}/recordings/start`, {
    method: 'POST',
    data: { call_id: callId },
  });
}

export async function stopRecording(recordingId: string) {
  return request<API.MessageResponse>(`${API}/recordings/stop`, {
    method: 'POST',
    data: { recording_id: recordingId },
  });
}

export async function getRecording(recordingId: string) {
  return request<API.DataResponse<API.VoiceRecording>>(`${API}/recordings/${recordingId}`);
}

export async function listTranscripts(params?: { call_id?: string; page?: number; page_size?: number }) {
  const payload = await request<unknown>(`${API}/transcripts`, { params });
  return normalizePaginatedResponse<API.VoiceTranscript>(payload);
}
