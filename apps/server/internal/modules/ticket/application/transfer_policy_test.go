package application

import "testing"

func TestBuildTransferAssignmentUpdateAssignsOpenTicket(t *testing.T) {
	updates, fromStatus, toStatus := BuildTransferAssignmentUpdate(7, "open")
	if fromStatus != "open" || toStatus != "assigned" {
		t.Fatalf("unexpected statuses: from=%s to=%s", fromStatus, toStatus)
	}
	if updates["agent_id"] != uint(7) {
		t.Fatalf("expected agent_id update, got %+v", updates)
	}
	if updates["status"] != "assigned" {
		t.Fatalf("expected assigned status update, got %+v", updates)
	}
}

func TestBuildTransferAssignmentUpdatePreservesNonOpenStatus(t *testing.T) {
	updates, fromStatus, toStatus := BuildTransferAssignmentUpdate(9, "pending")
	if fromStatus != "pending" || toStatus != "pending" {
		t.Fatalf("unexpected statuses: from=%s to=%s", fromStatus, toStatus)
	}
	if updates["agent_id"] != uint(9) {
		t.Fatalf("expected agent_id update, got %+v", updates)
	}
	if _, ok := updates["status"]; ok {
		t.Fatalf("did not expect status update, got %+v", updates)
	}
}

