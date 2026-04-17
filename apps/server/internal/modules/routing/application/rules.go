package application

import "fmt"

// NormalizeWaitingRecordQuery centralizes default waiting-queue query semantics
// so legacy runtime facades do not drift from the routing module.
func NormalizeWaitingRecordQuery(status string, limit int) (string, int) {
	if status == "" {
		status = "waiting"
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return status, limit
}

// BuildTransferMessage centralizes the default human-handoff notification text.
func BuildTransferMessage(reason, notes string) string {
	message := "您的会话已转接至人工客服"

	if reason != "" {
		message += fmt.Sprintf("。转接原因：%s", reason)
	}

	if notes != "" {
		message += fmt.Sprintf("。备注：%s", notes)
	}

	message += "。客服将很快为您提供帮助。"
	return message
}

// BuildWaitingMessage centralizes the default queueing notification text.
func BuildWaitingMessage() string {
	return "您的会话已加入人工客服等待队列，我们会尽快为您安排客服。请耐心等待。"
}

// BuildWaitingCancellationMessage centralizes the default queue cancellation text.
func BuildWaitingCancellationMessage(reason string) string {
	return fmt.Sprintf("已取消人工客服等待队列（原因：%s）", reason)
}

// BuildWaitingAlreadyQueuedSummary centralizes the idempotent queueing result text.
func BuildWaitingAlreadyQueuedSummary() string {
	return "会话已在等待队列中"
}

// BuildWaitingQueuedSummary centralizes the successful queueing result text.
func BuildWaitingQueuedSummary() string {
	return "会话已加入等待队列"
}

// BuildTransferAlreadyAssignedSummary centralizes the idempotent transfer result text.
func BuildTransferAlreadyAssignedSummary() string {
	return "会话已指派给目标客服"
}
