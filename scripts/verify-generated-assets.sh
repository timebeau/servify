#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="$ROOT_DIR/generated-assets.manifest"

if [ ! -f "$MANIFEST" ]; then
  echo "Missing generated assets manifest: $MANIFEST"
  exit 1
fi

while IFS= read -r asset; do
  if [ -z "$asset" ] || [[ "$asset" == \#* ]]; then
    continue
  fi

  if [ ! -f "$ROOT_DIR/$asset" ]; then
    echo "Missing generated asset: $asset"
    exit 1
  fi

  git -C "$ROOT_DIR" ls-files --error-unmatch "$asset" >/dev/null 2>&1 || {
    echo "Generated asset is not tracked by git: $asset"
    exit 1
  }
done < "$MANIFEST"

echo "Generated assets manifest verified."
