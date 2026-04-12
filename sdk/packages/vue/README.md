# @servify/vue

Vue surface for the Servify Web SDK.

## Usage

```ts
import { ServifyPlugin } from '@servify/vue';
```

Available composables:

- `useServify()`
- `useServifyReady()`
- `useChat()`
- `useAI()`
- `useTickets()`
- `useSatisfaction()`
- `useRemoteAssist()`

## Remote Assistance

`useRemoteAssist()` now exposes the same remote-assistance baseline as the React surface:

- `state`
- `isActive`
- `error`
- `remoteStream`
- `startRemoteAssist(options?)`
- `acceptRemoteAnswer(answer)`
- `addRemoteIce(candidate)`
- `endRemoteAssist()`
