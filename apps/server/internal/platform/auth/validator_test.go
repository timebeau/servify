package auth

import (
	"testing"
	"time"
)

func TestValidatorValidateToken(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 7,
		"iat":     now.Unix(),
		"exp":     now.Add(5 * time.Minute).Unix(),
	}, secret)

	payload, err := (Validator{
		Secret: secret,
		Now:    func() time.Time { return now },
	}).ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if got := payload["user_id"]; got != float64(7) {
		t.Fatalf("user_id = %v want 7", got)
	}
}

func TestValidatorValidateTokenExpired(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 7,
		"iat":     now.Add(-10 * time.Minute).Unix(),
		"exp":     now.Add(-5 * time.Minute).Unix(),
	}, secret)

	_, err := (Validator{
		Secret: secret,
		Now:    func() time.Time { return now },
	}).ValidateToken(token)
	if err == nil || err.Error() != "token time constraint failed: exp" {
		t.Fatalf("expected exp error, got %v", err)
	}
}
