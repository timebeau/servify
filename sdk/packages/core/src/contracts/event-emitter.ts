export type EventMap = Record<string, unknown[]>;

export interface EventEmitterLike<TEvents extends EventMap> {
  on<TKey extends keyof TEvents & string>(
    event: TKey,
    handler: (...args: TEvents[TKey]) => void
  ): void;
  off<TKey extends keyof TEvents & string>(
    event: TKey,
    handler?: (...args: TEvents[TKey]) => void
  ): void;
  emit<TKey extends keyof TEvents & string>(
    event: TKey,
    ...args: TEvents[TKey]
  ): void;
  removeAllListeners?(event?: keyof TEvents & string): void;
}
