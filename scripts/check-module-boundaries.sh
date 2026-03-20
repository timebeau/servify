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

require_pattern \
  "apps/server/internal/handlers/agent_handler.go" \
  'agentService\s+agentdelivery\.HandlerService' \
  "Agent handler must store modules/agent/delivery.HandlerService."
require_pattern \
  "apps/server/internal/handlers/agent_handler.go" \
  'func NewAgentHandler\(agentService agentdelivery\.HandlerService, logger \*logrus\.Logger\)' \
  "Agent handler constructor must accept modules/agent/delivery.HandlerService."
forbid_pattern \
  "apps/server/internal/handlers/agent_handler.go" \
  '\*services\.AgentService' \
  "Agent handler must not depend on concrete services.AgentService."

require_pattern \
  "apps/server/internal/handlers/ticket_handler.go" \
  'ticketService\s+ticketdelivery\.HandlerService' \
  "Ticket handler must store modules/ticket/delivery.HandlerService."
require_pattern \
  "apps/server/internal/handlers/ticket_handler.go" \
  'func NewTicketHandler\(ticketService ticketdelivery\.HandlerService, logger \*logrus\.Logger\)' \
  "Ticket handler constructor must accept modules/ticket/delivery.HandlerService."
forbid_pattern \
  "apps/server/internal/handlers/ticket_handler.go" \
  '\*services\.TicketService' \
  "Ticket handler must not depend on concrete services.TicketService."

require_pattern \
  "apps/server/internal/app/server/router.go" \
  'AgentHandlerService\s+agentdelivery\.HandlerService' \
  "Router dependencies must expose agentdelivery.HandlerService for agent handlers."
require_pattern \
  "apps/server/internal/app/server/router.go" \
  'TicketHandlerService\s+ticketdelivery\.HandlerService' \
  "Router dependencies must expose ticketdelivery.HandlerService for ticket handlers."
require_pattern \
  "apps/server/internal/app/server/runtime.go" \
  'AgentHandlerService\s+agentdelivery\.HandlerService' \
  "Runtime must keep agent handler wiring behind agentdelivery.HandlerService."
require_pattern \
  "apps/server/internal/app/server/runtime.go" \
  'TicketHandlerService\s+ticketdelivery\.HandlerService' \
  "Runtime must keep ticket handler wiring behind ticketdelivery.HandlerService."

require_pattern \
  "apps/server/internal/handlers/statistics_handler.go" \
  'statsService\s+analyticsdelivery\.HandlerService' \
  "Statistics handler must store modules/analytics/delivery.HandlerService."
require_pattern \
  "apps/server/internal/handlers/statistics_handler.go" \
  'func NewStatisticsHandler\(statsService analyticsdelivery\.HandlerService, logger \*logrus\.Logger\)' \
  "Statistics handler constructor must accept modules/analytics/delivery.HandlerService."
forbid_pattern \
  "apps/server/internal/handlers/statistics_handler.go" \
  '\*services\.StatisticsService' \
  "Statistics handler must not depend on concrete services.StatisticsService."
require_pattern \
  "apps/server/internal/app/server/router.go" \
  'StatisticsHandlerService\s+analyticsdelivery\.HandlerService' \
  "Router dependencies must expose analyticsdelivery.HandlerService for statistics handlers."
require_pattern \
  "apps/server/internal/app/server/runtime.go" \
  'StatisticsHandlerService\s+analyticsdelivery\.HandlerService' \
  "Runtime must keep statistics handler wiring behind analyticsdelivery.HandlerService."

require_pattern \
  "apps/server/internal/handlers/session_transfer_handler.go" \
  'transferService\s+routingdelivery\.HandlerService' \
  "Session transfer handler must store modules/routing/delivery.HandlerService."
require_pattern \
  "apps/server/internal/handlers/session_transfer_handler.go" \
  'func NewSessionTransferHandler\(transferService routingdelivery\.HandlerService, logger \*logrus\.Logger\)' \
  "Session transfer handler constructor must accept modules/routing/delivery.HandlerService."
forbid_pattern \
  "apps/server/internal/handlers/session_transfer_handler.go" \
  '\*services\.SessionTransferService' \
  "Session transfer handler must not depend on concrete services.SessionTransferService."
require_pattern \
  "apps/server/internal/app/server/router.go" \
  'TransferHandlerService\s+routingdelivery\.HandlerService' \
  "Router dependencies must expose routingdelivery.HandlerService for session transfer handlers."
require_pattern \
  "apps/server/internal/app/server/runtime.go" \
  'TransferHandlerService\s+routingdelivery\.HandlerService' \
  "Runtime must keep session transfer handler wiring behind routingdelivery.HandlerService."

if [[ "$has_error" -ne 0 ]]; then
  exit 1
fi

echo "Module boundary checks passed."
