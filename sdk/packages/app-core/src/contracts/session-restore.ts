export interface SessionSnapshot {
  sessionId: string;
  customerId?: string;
  lastMessageId?: string;
  transportCursor?: string;
  savedAt: string;
}

export interface SessionRestoreStrategy {
  capture(): Promise<SessionSnapshot | null>;
  restore(snapshot: SessionSnapshot): Promise<void>;
  clear(sessionId: string): Promise<void>;
}

export interface AppReconnectPolicy {
  maxAttempts: number;
  baseDelayMs: number;
  restoreSessionOnReconnect: boolean;
}
