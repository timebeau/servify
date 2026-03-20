export interface OfflineQueueEntry<TPayload = unknown> {
  id: string;
  type: string;
  payload: TPayload;
  createdAt: string;
  retryCount: number;
}

export interface OfflineQueueStore<TPayload = unknown> {
  enqueue(entry: OfflineQueueEntry<TPayload>): Promise<void>;
  peek(limit?: number): Promise<OfflineQueueEntry<TPayload>[]>;
  acknowledge(entryId: string): Promise<void>;
  clear(): Promise<void>;
}
