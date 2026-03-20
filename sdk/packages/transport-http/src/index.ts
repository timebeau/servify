export interface HttpRequestTransportConfig {
  baseUrl: string;
  timeoutMs?: number;
  defaultHeaders?: Record<string, string>;
}

export interface LongPollingTransportConfig extends HttpRequestTransportConfig {
  pollIntervalMs: number;
  maxInflightRequests?: number;
}

export interface HttpTransportReservation {
  request: HttpRequestTransportConfig;
  longPolling?: LongPollingTransportConfig;
}

export const TRANSPORT_HTTP_RESERVED = 'reserved-for-http-transport-package';
