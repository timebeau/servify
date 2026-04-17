package application

// BuildTransferAssignmentUpdate centralizes the minimal ticket state change
// needed when a conversation transfer assigns a ticket to another agent.
func BuildTransferAssignmentUpdate(targetAgentID uint, currentStatus string) (map[string]interface{}, string, string) {
	updates := map[string]interface{}{"agent_id": targetAgentID}
	fromStatus := currentStatus
	toStatus := fromStatus
	if fromStatus == "open" || fromStatus == "" {
		toStatus = "assigned"
		updates["status"] = toStatus
	}
	return updates, fromStatus, toStatus
}

