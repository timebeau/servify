package application

import "testing"

func TestNormalizeWaitingRecordQuery(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		limit      int
		wantStatus string
		wantLimit  int
	}{
		{name: "defaults", status: "", limit: 0, wantStatus: "waiting", wantLimit: 50},
		{name: "caps large limit", status: "transferred", limit: 500, wantStatus: "transferred", wantLimit: 50},
		{name: "keeps valid values", status: "cancelled", limit: 25, wantStatus: "cancelled", wantLimit: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotLimit := NormalizeWaitingRecordQuery(tt.status, tt.limit)
			if gotStatus != tt.wantStatus || gotLimit != tt.wantLimit {
				t.Fatalf("NormalizeWaitingRecordQuery(%q, %d) = (%q, %d), want (%q, %d)",
					tt.status, tt.limit, gotStatus, gotLimit, tt.wantStatus, tt.wantLimit)
			}
		})
	}
}

func TestBuildTransferMessage(t *testing.T) {
	got := BuildTransferMessage("用户请求人工", "VIP 客户")
	want := "您的会话已转接至人工客服。转接原因：用户请求人工。备注：VIP 客户。客服将很快为您提供帮助。"
	if got != want {
		t.Fatalf("BuildTransferMessage() = %q, want %q", got, want)
	}

	got = BuildTransferMessage("", "")
	want = "您的会话已转接至人工客服。客服将很快为您提供帮助。"
	if got != want {
		t.Fatalf("BuildTransferMessage() default = %q, want %q", got, want)
	}
}

func TestBuildWaitingMessage(t *testing.T) {
	got := BuildWaitingMessage()
	want := "您的会话已加入人工客服等待队列，我们会尽快为您安排客服。请耐心等待。"
	if got != want {
		t.Fatalf("BuildWaitingMessage() = %q, want %q", got, want)
	}
}

func TestBuildWaitingCancellationMessage(t *testing.T) {
	got := BuildWaitingCancellationMessage("用户离开")
	want := "已取消人工客服等待队列（原因：用户离开）"
	if got != want {
		t.Fatalf("BuildWaitingCancellationMessage() = %q, want %q", got, want)
	}
}

func TestBuildWaitingResultSummaries(t *testing.T) {
	if got := BuildWaitingAlreadyQueuedSummary(); got != "会话已在等待队列中" {
		t.Fatalf("BuildWaitingAlreadyQueuedSummary() = %q", got)
	}
	if got := BuildWaitingQueuedSummary(); got != "会话已加入等待队列" {
		t.Fatalf("BuildWaitingQueuedSummary() = %q", got)
	}
	if got := BuildTransferAlreadyAssignedSummary(); got != "会话已指派给目标客服" {
		t.Fatalf("BuildTransferAlreadyAssignedSummary() = %q", got)
	}
}
