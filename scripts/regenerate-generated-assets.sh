#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "==> Syncing demo SDK artifacts"
"$ROOT_DIR/scripts/sync-sdk-to-demo.sh"

echo "==> Regenerating API docs"
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g apps/server/cmd/server/main.go -o docs/generated/api

echo "==> Verifying generated assets manifest"
"$ROOT_DIR/scripts/verify-generated-assets.sh"

echo "Generated assets regenerated successfully."
