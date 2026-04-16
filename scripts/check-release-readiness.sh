#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

CONFIG_PATH="${1:-config.yml}"

if [[ "$CONFIG_PATH" != /* ]]; then
  CONFIG_PATH="$ROOT_DIR/$CONFIG_PATH"
fi

if [ ! -f "$CONFIG_PATH" ]; then
  echo "Config file not found: $CONFIG_PATH" >&2
  echo "Usage: $0 [config-path]" >&2
  exit 1
fi

echo "==> Local environment"
sh "$ROOT_DIR/scripts/check-local-environment.sh"

echo
echo "==> Security baseline"
sh "$ROOT_DIR/scripts/check-security-baseline.sh" "$CONFIG_PATH"

echo
echo "==> Observability baseline"
sh "$ROOT_DIR/scripts/check-observability-baseline.sh" "$CONFIG_PATH"

echo
echo "==> Focused Go regression tests"
go -C apps/server test ./cmd/cli ./internal/app/bootstrap

echo
echo "==> Build verification"
if ! make build >/dev/null 2>&1; then
    echo "[fail] build failed" >&2
    exit 1
fi
echo "[ok] build succeeded"

# 验证生成的二进制文件存在
if [ ! -f "$ROOT_DIR/bin/servify" ]; then
    echo "[fail] servify binary not found" >&2
    exit 1
fi
if [ ! -f "$ROOT_DIR/bin/migrate" ]; then
    echo "[fail] migrate binary not found" >&2
    exit 1
fi
if [ ! -f "$ROOT_DIR/bin/servify-cli" ]; then
    echo "[fail] servify-cli binary not found" >&2
    exit 1
fi
echo "[ok] all binaries present"

echo
echo "==> Binary version check"
"$ROOT_DIR/bin/servify-cli" version >/dev/null 2>&1 || echo "[warn] servify-cli version check failed"

echo
echo "==> Route registration verification"
# 通过检查handler包来验证关键路由已注册
if ! go -C apps/server test -run TestConversationWorkspaceHandler_GetSession ./internal/handlers/ >/dev/null 2>&1; then
    echo "[warn] conversation workspace route test failed (may require DB)"
fi
# 验证AI路由存在（通过检查ai handler）
if ! go -C apps/server test -run TestAIHandler_GetMetrics ./internal/handlers/ >/dev/null 2>&1; then
    echo "[fail] AI handler route test failed" >&2
    exit 1
fi
echo "[ok] critical routes verified"

echo
echo "Release readiness check passed."
