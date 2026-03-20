#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

has_error=0

require_pattern() {
  local file="$1"
  local pattern="$2"
  local message="$3"
  if ! rg -q --multiline "$pattern" "$file"; then
    echo "$message"
    has_error=1
  fi
}

forbid_pattern() {
  local file="$1"
  local pattern="$2"
  local message="$3"
  if rg -q --multiline "$pattern" "$file"; then
    echo "$message"
    has_error=1
  fi
}

check_handler_contract() {
  local label="$1"
  local file="$2"
  local field_name="$3"
  local contract_expr="$4"
  local constructor_pattern="$5"
  local forbidden_pattern="$6"

  require_pattern \
    "$file" \
    "${field_name}[[:space:]]+${contract_expr}" \
    "${label} handler must store ${contract_expr}."
  require_pattern \
    "$file" \
    "$constructor_pattern" \
    "${label} handler constructor must accept ${contract_expr}."
  forbid_pattern \
    "$file" \
    "$forbidden_pattern" \
    "${label} handler must not depend on its concrete legacy service."
}

check_runtime_contract() {
  local label="$1"
  local contract_field="$2"
  local contract_expr="$3"

  require_pattern \
    "apps/server/internal/app/server/router.go" \
    "${contract_field}[[:space:]]+${contract_expr}" \
    "Router dependencies must expose ${contract_expr} for ${label}."
  require_pattern \
    "apps/server/internal/app/server/runtime.go" \
    "${contract_field}[[:space:]]+${contract_expr}" \
    "Runtime must keep ${label} behind ${contract_expr}."
}

handler_specs=(
  "Agent|apps/server/internal/handlers/agent_handler.go|agentService|agentdelivery\\.HandlerService|func NewAgentHandler\\(agentService agentdelivery\\.HandlerService, logger \\*logrus\\.Logger\\)|\\*services\\.AgentService"
  "Ticket|apps/server/internal/handlers/ticket_handler.go|ticketService|ticketdelivery\\.HandlerService|func NewTicketHandler\\(ticketService ticketdelivery\\.HandlerService, logger \\*logrus\\.Logger\\)|\\*services\\.TicketService"
  "Statistics|apps/server/internal/handlers/statistics_handler.go|statsService|analyticsdelivery\\.HandlerService|func NewStatisticsHandler\\(statsService analyticsdelivery\\.HandlerService, logger \\*logrus\\.Logger\\)|\\*services\\.StatisticsService"
  "Session transfer|apps/server/internal/handlers/session_transfer_handler.go|transferService|routingdelivery\\.HandlerService|func NewSessionTransferHandler\\(transferService routingdelivery\\.HandlerService, logger \\*logrus\\.Logger\\)|\\*services\\.SessionTransferService"
  "AI|apps/server/internal/handlers/ai_handler.go|aiService|aidelivery\\.HandlerService|func NewAIHandler\\(aiService aidelivery\\.HandlerService\\)|aiService services\\.AIServiceInterface"
)

runtime_specs=(
  "agent handlers|AgentHandlerService|agentdelivery\\.HandlerService"
  "ticket handlers|TicketHandlerService|ticketdelivery\\.HandlerService"
  "statistics handlers|StatisticsHandlerService|analyticsdelivery\\.HandlerService"
  "session transfer handlers|TransferHandlerService|routingdelivery\\.HandlerService"
  "AI handlers|AIHandlerService|aidelivery\\.HandlerService"
)

for spec in "${handler_specs[@]}"; do
  IFS='|' read -r label file field_name contract_expr constructor_pattern forbidden_pattern <<<"$spec"
  check_handler_contract "$label" "$file" "$field_name" "$contract_expr" "$constructor_pattern" "$forbidden_pattern"
done

for spec in "${runtime_specs[@]}"; do
  IFS='|' read -r label contract_field contract_expr <<<"$spec"
  check_runtime_contract "$label" "$contract_field" "$contract_expr"
done

require_pattern \
  "apps/server/internal/handlers/health_enhanced.go" \
  'aiService aidelivery\.HandlerService' \
  "Enhanced health handler must store modules/ai/delivery.HandlerService."

require_pattern \
  "apps/server/internal/services/websocket.go" \
  'conversationWriter conversationdelivery\.WebSocketMessageWriter' \
  "WebSocket hub must store the public conversationdelivery.WebSocketMessageWriter contract."
require_pattern \
  "apps/server/internal/services/websocket.go" \
  'func \(h \*WebSocketHub\) SetConversationMessageWriter\(writer conversationdelivery\.WebSocketMessageWriter\)' \
  "WebSocket hub must accept conversationdelivery.WebSocketMessageWriter."
require_pattern \
  "apps/server/internal/app/server/runtime.go" \
  'SetConversationMessageWriter\(conversationdelivery\.NewWebSocketMessageAdapter\(conversationService\)\)' \
  "Main runtime must wire conversation websocket persistence through the module delivery adapter."
require_pattern \
  "apps/server/internal/app/server/realtime_runtime.go" \
  'SetConversationMessageWriter\(conversationdelivery\.NewWebSocketMessageAdapter\(conversationService\)\)' \
  "Realtime runtime must wire conversation websocket persistence through the module delivery adapter."

require_pattern \
  "apps/server/internal/app/server/ai_runtime.go" \
  'Service[[:space:]]+aidelivery\.HandlerService' \
  "AI assembly must expose the handler-facing AI contract."
require_pattern \
  "apps/server/internal/app/server/ai_runtime.go" \
  'RuntimeService[[:space:]]+services\.AIServiceInterface' \
  "AI assembly must keep a separate runtime AIServiceInterface for non-handler callers."

if [[ "$has_error" -ne 0 ]]; then
  exit 1
fi

echo "Module boundary checks passed."
