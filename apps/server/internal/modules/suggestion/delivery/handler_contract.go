package delivery

import (
	"context"

	suggestioncontract "servify/apps/server/internal/modules/suggestion/contract"
)

type HandlerService interface {
	Suggest(ctx context.Context, req *suggestioncontract.SuggestionRequest) (*suggestioncontract.SuggestionResponse, error)
}
