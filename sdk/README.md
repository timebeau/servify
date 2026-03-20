# Servify SDK Workspace

This workspace currently implements web SDK packages and reserves package slots for future SDK families.

## Current packages

- `packages/core`
- `packages/react`
- `packages/vue`
- `packages/vanilla`

## Reserved packages

- `packages/api-client`
- `packages/app-core`

Reserved packages now include stable design-time contracts so future SDK work can extend them without copying Web SDK internals.

## Target structure

- `core`
- `transport-http`
- `transport-websocket`
- `web-vanilla`
- `web-react`
- `web-vue`
- `api-client`
- `app-core`
- `transport-http`
- `transport-websocket`

Current implementation remains web-only. Reserved packages are placeholders for future expansion and should not contain production logic yet.
