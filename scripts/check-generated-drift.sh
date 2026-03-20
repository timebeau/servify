#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -lt 3 ]; then
  echo "Usage: $0 <label> <regenerate-command> <path> [<path> ...]"
  exit 2
fi

LABEL="$1"
REGENERATE_COMMAND="$2"
shift 2

if git diff --quiet -- "$@"; then
  echo "$LABEL drift check passed."
  exit 0
fi

echo "$LABEL drift detected."
echo
echo "Affected paths:"
for path in "$@"; do
  echo "  - $path"
done
echo
echo "Diff summary:"
git diff --stat -- "$@"
echo
echo "Regenerate with:"
echo "  $REGENERATE_COMMAND"
echo
echo "Then review and commit the updated files."
exit 1
