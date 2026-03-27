package auth

import (
	"encoding/json"
	"errors"
	"time"
)

type TokenPolicy func(payload map[string]interface{}, claims Claims, now time.Time) error

func ComposeTokenPolicies(policies ...TokenPolicy) TokenPolicy {
	return func(payload map[string]interface{}, claims Claims, now time.Time) error {
		for _, policy := range policies {
			if policy == nil {
				continue
			}
			if err := policy(payload, claims, now); err != nil {
				return err
			}
		}
		return nil
	}
}

func RejectIssuedBefore(minIssuedAt int64) TokenPolicy {
	return func(payload map[string]interface{}, claims Claims, now time.Time) error {
		if minIssuedAt <= 0 {
			return nil
		}
		issuedAt, ok := int64Claim(payload, "iat")
		if !ok {
			return errors.New("token missing iat required by revocation policy")
		}
		if issuedAt < minIssuedAt {
			return errors.New("token has been revoked by issued-at policy")
		}
		return nil
	}
}

func RequireMinimumTokenVersion(minVersion int64) TokenPolicy {
	return func(payload map[string]interface{}, claims Claims, now time.Time) error {
		if minVersion <= 0 {
			return nil
		}
		version, ok := int64Claim(payload, "token_version", "ver")
		if !ok {
			return errors.New("token missing token_version required by session policy")
		}
		if version < minVersion {
			return errors.New("token has been revoked by version policy")
		}
		return nil
	}
}

func int64Claim(payload map[string]interface{}, keys ...string) (int64, bool) {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return int64(typed), true
		case json.Number:
			n, err := typed.Int64()
			if err == nil {
				return n, true
			}
		}
	}
	return 0, false
}
