package application

import (
	"fmt"
	"strings"
)

const (
	defaultMaxQueryLength    = 4000
	defaultMaxOutputLength   = 2000
	guardrailBlockedResponse = "请求包含受限内容，无法处理。"
)

// Guardrails applies basic input and output protections for AI requests.
type Guardrails struct {
	maxQueryLength  int
	maxOutputLength int
	blockedTerms    []string
}

func NewGuardrails() *Guardrails {
	return &Guardrails{
		maxQueryLength:  defaultMaxQueryLength,
		maxOutputLength: defaultMaxOutputLength,
		blockedTerms: []string{
			"DROP TABLE",
			"rm -rf",
		},
	}
}

func (g *Guardrails) ValidateInput(req AIRequest) error {
	query := strings.TrimSpace(req.Query)
	if query == "" && len(req.Messages) == 0 {
		return fmt.Errorf("empty ai request")
	}
	if len(query) > g.maxQueryLength {
		return fmt.Errorf("query exceeds max length")
	}

	upper := strings.ToUpper(query)
	for _, term := range g.blockedTerms {
		if strings.Contains(upper, term) {
			return fmt.Errorf("request blocked by guardrails")
		}
	}
	return nil
}

func (g *Guardrails) SanitizeOutput(content string) (string, bool) {
	if len(content) <= g.maxOutputLength {
		return content, false
	}
	return content[:g.maxOutputLength], true
}
