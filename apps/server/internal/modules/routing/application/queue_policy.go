package application

import "strings"

// IsActiveWaitingStatus reports whether a queue record should still participate
// in waiting-queue processing and duplicate-handoff suppression.
func IsActiveWaitingStatus(status string) bool {
	return strings.EqualFold(strings.TrimSpace(status), "waiting")
}
