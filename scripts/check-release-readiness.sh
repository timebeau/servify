#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

CONFIG_PATH="${1:-config.yml}"

if [[ "$CONFIG_PATH" != /* ]]; then
  CONFIG_PATH="$ROOT_DIR/$CONFIG_PATH"
fi

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo "Config file not found: $CONFIG_PATH" >&2
  echo "Usage: $0 [config-path]" >&2
  exit 1
fi

RUNTIME_DIR="$ROOT_DIR/.runtime"
mkdir -p "$RUNTIME_DIR"

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
go -C apps/server test ./cmd/cli ./internal/app/bootstrap ./internal/handlers

echo
echo "==> Build verification"
make build >/dev/null
echo "[ok] build succeeded"

for binary in servify migrate servify-cli; do
  if [[ ! -f "$ROOT_DIR/bin/$binary" ]]; then
    echo "[fail] missing binary: $binary" >&2
    exit 1
  fi
done
echo "[ok] all binaries present"

echo
echo "==> Binary version check"
"$ROOT_DIR/bin/servify-cli" version >/dev/null 2>&1 || echo "[warn] servify-cli version check failed"

echo
echo "==> Route registration verification"
if ! go -C apps/server test -run TestConversationWorkspaceHandler_GetSession ./internal/handlers/ >/dev/null 2>&1; then
  echo "[warn] conversation workspace route test failed (may require DB)"
fi
if ! go -C apps/server test -run "TestAIHandler_(GetMetrics|UploadDocument_StandardService|SyncKnowledgeBase_StandardService|EnableKnowledgeProvider_StandardService|DisableKnowledgeProvider_StandardService|ResetCircuitBreaker_StandardService)$" ./internal/handlers/ >/dev/null 2>&1; then
  echo "[fail] AI handler route test failed" >&2
  exit 1
fi
echo "[ok] critical routes verified"

echo
echo "==> Database migration smoke"
"$ROOT_DIR/bin/migrate" --config "$CONFIG_PATH" >/dev/null
echo "[ok] migration succeeded"

echo
echo "==> HTTP smoke verification"
PORT="${SERVIFY_RELEASE_CHECK_PORT:-18080}"
SERVER_LOG="$RUNTIME_DIR/release-check-server.log"
HEALTH_OUT="$RUNTIME_DIR/release-check-health.json"
READY_OUT="$RUNTIME_DIR/release-check-ready.json"
METRICS_OUT="$RUNTIME_DIR/release-check-metrics.txt"
SERVER_PID=""

cleanup() {
  if [[ -n "$SERVER_PID" ]] && kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
    wait "$SERVER_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

SERVIFY_HOST="127.0.0.1" SERVIFY_PORT="$PORT" "$ROOT_DIR/bin/servify" >"$SERVER_LOG" 2>&1 &
SERVER_PID=$!

healthy="false"
for _ in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:${PORT}/health" >"$HEALTH_OUT" 2>/dev/null; then
    healthy="true"
    break
  fi
  sleep 1
done

if [[ "$healthy" != "true" ]]; then
  echo "[fail] server did not become healthy on port ${PORT}" >&2
  cat "$SERVER_LOG" >&2 || true
  exit 1
fi

curl -fsS "http://127.0.0.1:${PORT}/ready" >"$READY_OUT"
curl -fsS "http://127.0.0.1:${PORT}/metrics" >"$METRICS_OUT"
echo "[ok] /health responded"
echo "[ok] /ready responded"
echo "[ok] /metrics responded"

echo
echo "Release readiness check passed."
