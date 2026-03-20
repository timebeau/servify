#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Running Go smoke: voice protocol runtime"
go test -tags=integration ./apps/server/internal/handlers -run 'TestVoiceHandler.*Integration' -count=1

echo
echo "Running Go smoke: AI fallback wiring"
go test -tags=integration ./apps/server/internal/services -run 'TestOrchestratedEnhancedAIServiceFallback.*Integration' -count=1

echo
echo "Running SDK example smoke tests"
npm -C sdk run test:examples
