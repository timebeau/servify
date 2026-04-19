#!/usr/bin/env bash

set -euo pipefail

FROM_REF="${1:-}"
TO_REF="${2:-HEAD}"

if [ -z "$FROM_REF" ]; then
  last_tag="$(git describe --tags --abbrev=0 2>/dev/null || true)"
  if [ -n "$last_tag" ]; then
    FROM_REF="$last_tag"
  else
    FROM_REF="$(git rev-list --max-parents=0 HEAD)"
  fi
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

echo "# Changelog Draft"
echo
echo "Range: \`$FROM_REF..$TO_REF\`"
echo
echo "## Commits"
git log --no-merges --pretty=format:'- %h %s (%an)' "$FROM_REF..$TO_REF"
echo
