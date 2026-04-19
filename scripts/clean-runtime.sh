#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

paths=(
  ".runtime"
  "uploads"
  "apps/server/uploads"
  "internal/handlers/uploads"
)

for path in "${paths[@]}"; do
  if [ -e "$path" ]; then
    rm -rf "$path"
    echo "Removed: $path"
  else
    echo "Skipped missing path: $path"
  fi
done

echo "Runtime cleanup completed."
