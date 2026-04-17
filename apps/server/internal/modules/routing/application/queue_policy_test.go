package application

import "testing"

func TestIsActiveWaitingStatus(t *testing.T) {
	if !IsActiveWaitingStatus("waiting") {
		t.Fatal("expected waiting to be active")
	}
	if !IsActiveWaitingStatus(" WAITING ") {
		t.Fatal("expected trimmed waiting to be active")
	}
	if IsActiveWaitingStatus("cancelled") {
		t.Fatal("did not expect cancelled to be active")
	}
	if IsActiveWaitingStatus("transferred") {
		t.Fatal("did not expect transferred to be active")
	}
}
