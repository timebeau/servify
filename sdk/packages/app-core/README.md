# @servify/app-core

保留包，用于未来 App / Mobile SDK 的共享基础层。

预留场景：

- mobile app session
- push token integration
- reconnect and offline queue

当前不实现具体 runtime，但已经预留以下移动端核心 contract：

- `OfflineQueueStore`
- `PushTokenRegistrar`
- `SessionRestoreStrategy`
- `MobileStorageAdapter`

这些 contract 用于后续 iOS / Android / React Native SDK 收口，避免直接复制 Web SDK 结构。
