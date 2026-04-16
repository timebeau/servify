#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "⚠️  scripts/sync-sdk-to-admin.sh 已废弃，请改用 scripts/sync-sdk-to-demo.sh"
exec "$ROOT_DIR/scripts/sync-sdk-to-demo.sh" "$@"
