export interface WebSocketTransportConfig {
  url: string;
  protocols?: string[];
  heartbeatIntervalMs?: number;
}

export interface ServerSentEventsTransportConfig {
  url: string;
  withCredentials?: boolean;
  eventTypes?: string[];
}

export interface WebhookCallbackTransportConfig {
  callbackUrl: string;
  secretHeader?: string;
  challengePath?: string;
}

export interface WebSocketTransportReservation {
  websocket: WebSocketTransportConfig;
  sse?: ServerSentEventsTransportConfig;
  webhookCallback?: WebhookCallbackTransportConfig;
}

export const TRANSPORT_WEBSOCKET_RESERVED = 'reserved-for-websocket-transport-package';
