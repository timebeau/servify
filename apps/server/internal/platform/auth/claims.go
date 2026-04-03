package auth

import (
	"encoding/json"
	"strings"
)

// Claims is the normalized auth context derived from a validated JWT payload.
type Claims struct {
	Values        map[string]interface{}
	UserID        uint
	HasUserID     bool
	UserIDRaw     interface{}
	SessionID     string
	Roles         []string
	Permissions   []string
	PrincipalKind string
	TenantID      string
	WorkspaceID   string
}

func extractClaims(payload map[string]interface{}, resolver Resolver) Claims {
	claims := Claims{
		Values: payload,
		Roles:  normalizeStringList(payload["roles"]),
	}
	claims.PrincipalKind = derivePrincipalKind(payload, claims.Roles)
	claims.TenantID = normalizeOptionalString(firstValue(payload["tenant_id"], payload["tenant"], payload["tid"]))
	claims.WorkspaceID = normalizeOptionalString(firstValue(payload["workspace_id"], payload["workspace"], payload["wid"]))
	claims.SessionID = normalizeOptionalString(firstValue(payload["session_id"], payload["sid"]))

	if uid, ok := firstNonNil(payload["user_id"], payload["sub"]); ok {
		switch t := uid.(type) {
		case float64:
			claims.UserID = uint(t)
			claims.HasUserID = true
		case json.Number:
			if n, err := t.Int64(); err == nil {
				claims.UserID = uint(n)
				claims.HasUserID = true
			} else {
				claims.UserIDRaw = uid
			}
		default:
			claims.UserIDRaw = uid
		}
	}

	explicit := normalizeStringList(firstValue(payload["perms"], payload["permissions"]))
	claims.Permissions = resolver.ExpandPermissions(claims.Roles, explicit)
	return claims
}

func firstValue(vals ...interface{}) interface{} {
	for _, v := range vals {
		if v != nil {
			return v
		}
	}
	return nil
}

func firstNonNil(vals ...interface{}) (interface{}, bool) {
	for _, v := range vals {
		if v != nil {
			return v, true
		}
	}
	return nil, false
}

func normalizeStringList(v interface{}) []string {
	switch t := v.(type) {
	case nil:
		return nil
	case []string:
		out := make([]string, 0, len(t))
		for _, s := range t {
			if s = strings.TrimSpace(s); s != "" {
				out = append(out, s)
			}
		}
		return out
	case []interface{}:
		var out []string
		for _, it := range t {
			if s, ok := it.(string); ok {
				if s = strings.TrimSpace(s); s != "" {
					out = append(out, s)
				}
			}
		}
		return out
	case string:
		if strings.TrimSpace(t) == "" {
			return nil
		}
		parts := strings.Split(t, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func dedupeStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func normalizeOptionalString(v interface{}) string {
	s, _ := v.(string)
	return strings.TrimSpace(s)
}
