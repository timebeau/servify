package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func createTestHS256JWT(t *testing.T, payload map[string]interface{}, secret string) string {
	t.Helper()

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
	h := enc(headerJSON)
	p := enc(payloadJSON)
	signing := h + "." + p

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signing))
	sig := enc(mac.Sum(nil))
	return signing + "." + sig
}
