# @servify/vanilla

Vanilla JavaScript surface for the Servify Web SDK.

## Usage

```html
<script src="../../packages/vanilla/dist/index.js"></script>
```

Examples live in `sdk/examples/vanilla`.

## Remote Assistance

The vanilla surface now forwards the core Web SDK remote-assistance API:

- `startRemoteAssist(options?)`
- `acceptRemoteAnswer(answer)`
- `addRemoteIce(candidate)`
- `endRemoteAssist()`

Forwarded events:

- `webrtc:offer`
- `webrtc:answer`
- `webrtc:candidate`
- `webrtc:track`
- `webrtc:state`
