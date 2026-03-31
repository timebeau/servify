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
echo "Release readiness check passed."
