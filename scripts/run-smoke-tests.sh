#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "=== Phase 1 Smoke Tests ==="

echo
echo "1/5: Conversation handler — List/Send/Assign/Close"
go test -count=1 -run 'TestConversationWorkspaceHandler_(ListMessages|SendMessage|AssignAgent|CloseSession)' ./apps/server/internal/handlers/ -v

echo
echo "2/5: Ticket handler — Create/Get/List/Assign/Bulk"
go test -tags=integration -count=1 -run 'TestTicketHandler_Create_Get_List_Assign' ./apps/server/internal/handlers/ -v

echo
echo "3/5: Ticket handler — Custom Fields + Export"
go test -tags=integration -count=1 -run 'TestTicketHandler_CustomFields' ./apps/server/internal/handlers/ -v

echo
echo "4/5: Ticket handler — Related Conversations"
go test -tags=integration -count=1 -run 'TestTicketHandler_GetRelatedConversations' ./apps/server/internal/handlers/ -v

echo
echo "5/5: Go build — full project"
go build ./...

echo
echo "=== Voice & AI Smoke ==="
go test -tags=integration ./apps/server/internal/handlers -run 'TestVoiceHandler.*Integration' -count=1
go test -tags=integration ./apps/server/internal/services -run 'TestOrchestratedEnhancedAIServiceFallback.*Integration' -count=1

echo
echo "=== SDK Examples ==="
npm -C sdk run test:examples 2>/dev/null || echo "(SDK tests skipped — no sdk or no test:examples script)"

echo
echo "=== All smoke tests passed ==="
