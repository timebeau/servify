package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Validator verifies HS256 JWTs.
type Validator struct {
	Secret string
	Now    func() time.Time
}

// ValidateToken verifies a JWT and returns its payload.
func (v Validator) ValidateToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	headerB64, payloadB64, sigB64 := parts[0], parts[1], parts[2]

	headerJSON, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return nil, errors.New("invalid header encoding")
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, errors.New("invalid header json")
	}
	if alg, _ := header["alg"].(string); alg != "" && alg != "HS256" {
		return nil, errors.New("unsupported alg")
	}

	mac := hmac.New(sha256.New, []byte(v.Secret))
	mac.Write([]byte(headerB64 + "." + payloadB64))
	expected := mac.Sum(nil)
	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}
	if !hmac.Equal(sig, expected) {
		return nil, errors.New("invalid signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, errors.New("invalid payload json")
	}

	now := time.Now()
	if v.Now != nil {
		now = v.Now()
	}
	if err := validateTimeClaims(payload, now.Unix()); err != nil {
		return nil, err
	}

	return payload, nil
}

func validateTimeClaims(payload map[string]interface{}, nowSec int64) error {
	checkTime := func(key string, cmp func(int64) bool) error {
		if v, ok := payload[key]; ok {
			switch t := v.(type) {
			case float64:
				if !cmp(int64(t)) {
					return errors.New("token time constraint failed: " + key)
				}
			case json.Number:
				sec, _ := t.Int64()
				if !cmp(sec) {
					return errors.New("token time constraint failed: " + key)
				}
			}
		}
		return nil
	}

	if err := checkTime("nbf", func(sec int64) bool { return nowSec >= sec }); err != nil {
		return err
	}
	if err := checkTime("iat", func(sec int64) bool { return nowSec >= sec }); err != nil {
		return err
	}
	if err := checkTime("exp", func(sec int64) bool { return nowSec < sec }); err != nil {
		return err
	}
	return nil
}
