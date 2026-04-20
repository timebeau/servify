import { request } from '@/lib/request';

const API = '/api/voice';

type VoiceProtocolsResponse = {
  success: boolean;
  data: API.PaginatedResponse<API.VoiceProtocol>;
};

type VoiceTranscriptsResponse<T = API.VoiceTranscript> = {
  success: boolean;
  data: T[];
  total: number;
  page: number;
  page_size: number;
};

export async function listProtocols(params?: { page?: number; page_size?: number }) {
  return request<VoiceProtocolsResponse>(`${API}/protocols`, { params });
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
  return request<VoiceTranscriptsResponse>(`${API}/transcripts`, { params });
}
