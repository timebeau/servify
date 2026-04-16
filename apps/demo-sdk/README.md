# Servify Web SDK

`apps/demo-sdk` 保存浏览器可直接引用的 SDK 产物，以及一个轻量聊天挂件示例。

## 目录说明

- `servify-sdk.esm.js`：由 `sdk/packages/vanilla/dist/index.esm.js` 同步而来
- `servify-sdk.umd.js`：由 `sdk/packages/vanilla/dist/index.js` 同步而来，挂载到 `window.Servify`
- `index.d.ts`：由 `sdk/packages/vanilla/dist/index.d.ts` 同步而来
- `widget.js` / `widget.css`：演示用轻量挂件壳层，保留在仓库中，不由 SDK 构建自动产出

## 当前协议口径

- 默认实时通道：`/api/v1/ws`
- 聊天主链路：WebSocket-first，不再依赖旧的 `/api/sessions`、`/api/messages`
- AI：`/api/v1/ai/query`、`/api/v1/ai/status`
- 上传：`/api/v1/upload`
- 满意度：`/api/satisfactions`
- 会话查询与关闭：`/api/omni/sessions/*`
- 当前服务端未公开支持：
  - 队列 REST API
  - WebRTC call REST API
  - 旧式 REST 会话创建

## 浏览器直接使用

```html
<script src="/demo-sdk/servify-sdk.umd.js"></script>
<script>
  const client = new Servify({
    apiUrl: 'http://localhost:8080',
    wsUrl: 'ws://localhost:8080/api/v1/ws',
    customerName: 'Demo User',
    customerEmail: 'demo@example.com',
    debug: true,
  });

  await client.init();
  await client.startChat({ message: '你好，我需要帮助。' });
  await client.sendMessage('请帮我排查一下当前问题。');
  client.on('webrtc:state', (state) => console.log('remote assist state:', state));
</script>
```

## 远程协助

当前服务端采用会话级 WebRTC 信令模型，SDK 会把服务端 `webrtc-state-change` 归一到 `webrtc:state` 事件。

```ts
const peer = await client.startRemoteAssist({
  captureScreen: true,
  audio: false,
});

client.on('webrtc:answer', (answer) => client.acceptRemoteAnswer(answer));
client.on('webrtc:candidate', (candidate) => client.addRemoteIce(candidate));
client.on('webrtc:track', (event) => {
  const [stream] = event.streams;
  console.log('remote media stream:', stream);
});

await client.endRemoteAssist();
```

当前不承诺：

- 完整 co-browsing UI
- 专门坐席协助工作台
- 双端标准化演示系统

## 示例入口

- React：`sdk/examples/react`
- Vue：`sdk/examples/vue`
- Vanilla：`sdk/examples/vanilla`

## 重建与同步

```sh
npm -C sdk run build
bash ./scripts/sync-sdk-to-demo.sh
```

也可以直接执行：

```sh
make demo-sync-sdk
```
