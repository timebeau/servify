# Servify Web SDK

当前目录保存同步后的浏览器 SDK 产物，供 `apps/demo-sdk` 或外部静态页面直接接入。

文件：
- `servify-sdk.esm.js`：ES Module 构建
- `servify-sdk.umd.js`：UMD 构建，挂载到 `window.Servify`
- `index.d.ts`：类型声明

当前能力范围：
- Web 会话创建与消息发送
- AI 问答、文件上传、满意度提交
- 会话级 WebSocket 实时事件
- 远程协助基础链路：`startRemoteAssist`、`acceptRemoteAnswer`、`addRemoteIce`、`endRemoteAssist`
- WebRTC 事件：`webrtc:offer`、`webrtc:answer`、`webrtc:candidate`、`webrtc:track`、`webrtc:state`

## 浏览器直接使用

```html
<script src="/sdk/servify-sdk.umd.js"></script>
<script>
  const client = new Servify({
    apiUrl: 'http://localhost:8080',
    customerName: 'Demo User',
    customerEmail: 'demo@example.com',
    debug: true
  });

  await client.init();
  await client.startChat({ message: '你好，我需要帮助' });
  await client.sendMessage('请协助我排查问题');

  client.on('webrtc:state', (state) => {
    console.log('remote assist state:', state);
  });
</script>
```

## 远程协助

当前服务端是会话级信令模型，SDK 会把服务端 `webrtc-state-change` 也归一到 `webrtc:state` 事件。

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

- React 示例：`sdk/examples/react`
- Vue 示例：`sdk/examples/vue`
- Vanilla 示例：`sdk/examples/vanilla`

如果需要重新同步当前目录产物，执行：

```sh
sh scripts/sync-sdk-to-admin.sh
```
