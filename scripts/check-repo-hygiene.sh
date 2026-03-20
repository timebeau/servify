#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

tracked_forbidden=(
  "server"
  "server.exe"
)

tracked_forbidden_prefixes=(
  "uploads/"
  "apps/server/uploads/"
  "internal/handlers/uploads/"
  ".runtime/"
)

has_error=0

for path in "${tracked_forbidden[@]}"; do
  if git ls-files --error-unmatch -- "$path" >/dev/null 2>&1; then
    echo "Tracked runtime/build artifact is not allowed: $path"
    has_error=1
  fi
done

while IFS= read -r file; do
  for prefix in "${tracked_forbidden_prefixes[@]}"; do
    if [[ "$file" == "$prefix"* ]]; then
      echo "Tracked runtime/build artifact under forbidden path: $file"
      has_error=1
    fi
  done
done < <(git ls-files)

if [[ "$has_error" -ne 0 ]]; then
  exit 1
fi

echo "Repository hygiene checks passed."
