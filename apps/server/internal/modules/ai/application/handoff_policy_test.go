package application

import (
	"testing"

	"servify/apps/server/internal/models"
)

func TestShouldTransferToHuman(t *testing.T) {
	if !ShouldTransferToHuman("请帮我转人工客服", nil) {
		t.Fatal("expected handoff for human keyword")
	}
	if !ShouldTransferToHuman("我要投诉", nil) {
		t.Fatal("expected handoff for complaint keyword")
	}
	if !ShouldTransferToHuman("普通问题", make([]models.Message, 6)) {
		t.Fatal("expected handoff for long session history")
	}
	if ShouldTransferToHuman("普通问题", nil) {
		t.Fatal("did not expect handoff for normal query")
	}
}

func TestSimpleSessionSummary(t *testing.T) {
	if got := SimpleSessionSummary(nil); got != "空会话" {
		t.Fatalf("SimpleSessionSummary(nil) = %q", got)
	}

	got := SimpleSessionSummary([]models.Message{{Sender: "user", Content: " 你好，这是一个测试消息 "}})
	if got != "user: 你好，这是一个测试消息" {
		t.Fatalf("unexpected summary: %q", got)
	}

	long := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	got = SimpleSessionSummary([]models.Message{{Content: long}})
	if len(got) == 0 || len(got) > len("unknown: ")+80 {
		t.Fatalf("unexpected truncated summary length: %d", len(got))
	}
}

func TestBuildTransferSessionSummaryFallback(t *testing.T) {
	got := BuildTransferSessionSummaryFallback("alice", 7, []models.Message{{}, {}})
	if got != "用户alice的简短会话，共2条消息" {
		t.Fatalf("unexpected fallback summary: %q", got)
	}

	got = BuildTransferSessionSummaryFallback("", 42, nil)
	if got != "用户ID=42的简短会话，共0条消息" {
		t.Fatalf("unexpected id fallback summary: %q", got)
	}
}

func TestBuildSessionSummaryUnavailable(t *testing.T) {
	if got := BuildSessionSummaryUnavailable(); got != "无法生成会话摘要" {
		t.Fatalf("BuildSessionSummaryUnavailable() = %q", got)
	}
}
