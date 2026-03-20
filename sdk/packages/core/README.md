# @servify/core

Shared Web SDK runtime and contract package.

## Public API

- `createServifySDK`
- `createWebServifySDK`
- `WebSocketManager`
- exported contracts from `src/index.ts`

## Usage

```ts
import { createWebServifySDK } from '@servify/core';

const sdk = createWebServifySDK({
  apiUrl: 'http://localhost:8080',
  autoConnect: false,
});
```
