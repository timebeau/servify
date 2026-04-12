# @servify/core

Shared Web SDK runtime and contract package.

## Public API

- `createServifySDK`
- `createWebServifySDK`
- `WebSocketManager`
- `ServifySDK.startRemoteAssist`
- `ServifySDK.acceptRemoteAnswer`
- `ServifySDK.addRemoteIce`
- `ServifySDK.endRemoteAssist`
- exported contracts from `src/index.ts`

## Usage

```ts
import { createWebServifySDK } from '@servify/core';

const sdk = createWebServifySDK({
  apiUrl: 'http://localhost:8080',
  autoConnect: false,
});
```

## Remote Assistance

`@servify/core` now exposes the Web SDK remote-assistance baseline:

- capability flag: `remote_assist@1`
- WebSocket signaling events: `webrtc:offer`, `webrtc:answer`, `webrtc:candidate`, `webrtc:state`
- `webrtc:state` now also reflects server-pushed `webrtc-state-change` runtime updates
- runtime methods:
  - `startRemoteAssist({ captureScreen?: boolean, audio?: boolean, iceServers?: RTCIceServer[] })`
  - `acceptRemoteAnswer(answer)`
  - `addRemoteIce(candidate)`
  - `endRemoteAssist()`

Current scope is session-level WebRTC signaling and optional screen capture. It is not a full co-browsing UI.
