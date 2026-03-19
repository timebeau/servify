import type { CapabilitySet } from './capability';
import type { AuthProvider } from './auth-provider';
import type { EventEmitterLike, EventMap } from './event-emitter';
import type { Transport } from './transport';

export interface SessionIdentity {
  customerId?: string;
  conversationId?: string;
  agentId?: string;
}

export interface ClientSession<
  State = Record<string, unknown>,
  Events extends EventMap = Record<string, unknown[]>
> {
  readonly id: string;
  readonly transport: Transport;
  readonly authProvider?: AuthProvider;
  readonly capabilities: CapabilitySet;
  readonly events: EventEmitterLike<Events>;
  getIdentity(): SessionIdentity;
  getState(): State;
  updateState(patch: Partial<State>): State;
}
