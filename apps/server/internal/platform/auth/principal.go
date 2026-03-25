package auth

import "strings"

const (
	PrincipalUnknown = "unknown"
	PrincipalEndUser = "end_user"
	PrincipalAgent   = "agent"
	PrincipalAdmin   = "admin"
	PrincipalService = "service"
)

func normalizePrincipalKind(v interface{}) string {
	s, _ := v.(string)
	switch strings.ToLower(strings.TrimSpace(s)) {
	case PrincipalEndUser:
		return PrincipalEndUser
	case PrincipalAgent:
		return PrincipalAgent
	case PrincipalAdmin:
		return PrincipalAdmin
	case PrincipalService:
		return PrincipalService
	default:
		return ""
	}
}

func derivePrincipalKind(payload map[string]interface{}, roles []string) string {
	if kind := normalizePrincipalKind(firstValue(
		payload["principal_kind"],
		payload["principal_type"],
		payload["subject_type"],
		payload["token_type"],
	)); kind != "" {
		return kind
	}

	for _, role := range roles {
		switch strings.ToLower(strings.TrimSpace(role)) {
		case PrincipalAdmin, "super_admin":
			return PrincipalAdmin
		case PrincipalAgent:
			return PrincipalAgent
		}
	}

	if _, ok := firstNonNil(payload["user_id"], payload["sub"]); ok {
		return PrincipalEndUser
	}

	return PrincipalUnknown
}
