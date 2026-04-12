# @servify/react

React surface for the Servify Web SDK.

## Usage

```tsx
import { ServifyProvider } from '@servify/react';
```

Examples live in `sdk/examples/react`.

## Remote Assistance

The React surface now provides `useRemoteAssist()` on top of `@servify/core`.

Returned API:

- `state`
- `isActive`
- `remoteStream`
- `startRemoteAssist(options?)`
- `acceptRemoteAnswer(answer)`
- `addRemoteIce(candidate)`
- `endRemoteAssist()`
