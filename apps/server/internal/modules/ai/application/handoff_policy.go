package application

import (
	"fmt"
	"strings"

	"servify/apps/server/internal/models"
)

var humanHandoffKeywords = []string{
	"人工",
	"客服",
	"转人工",
	"manual",
	"human",
	"agent",
}

// ShouldTransferToHuman centralizes the default handoff heuristic used by
// legacy AI facades while the AI module becomes the primary business entry.
func ShouldTransferToHuman(query string, sessionHistory []models.Message) bool {
	query = strings.ToLower(query)

	for _, keyword := range humanHandoffKeywords {
		if strings.Contains(query, keyword) {
			return true
		}
	}

	if strings.Contains(query, "投诉") || strings.Contains(query, "complaint") {
		return true
	}

	return len(sessionHistory) > 5
}

// SimpleSessionSummary provides a deterministic summary fallback for runtimes
// that cannot or should not call an external model.
func SimpleSessionSummary(messages []models.Message) string {
	if len(messages) == 0 {
		return "空会话"
	}

	last := messages[len(messages)-1]
	content := strings.TrimSpace(last.Content)
	if content == "" {
		content = "无内容"
	}
	if len(content) > 80 {
		content = content[:80]
	}

	sender := strings.TrimSpace(last.Sender)
	if sender == "" {
		sender = "unknown"
	}

	return fmt.Sprintf("%s: %s", sender, content)
}

// BuildTransferSessionSummaryFallback provides a deterministic fallback summary
// for human handoff when the conversation is too short or AI summarization is unavailable.
func BuildTransferSessionSummaryFallback(customerLabel string, customerID uint, messages []models.Message) string {
	label := strings.TrimSpace(customerLabel)
	if label == "" {
		label = fmt.Sprintf("ID=%d", customerID)
	}
	return fmt.Sprintf("用户%s的简短会话，共%d条消息", label, len(messages))
}

// BuildSessionSummaryUnavailable provides the default summary when summary generation fails.
func BuildSessionSummaryUnavailable() string {
	return "无法生成会话摘要"
}
