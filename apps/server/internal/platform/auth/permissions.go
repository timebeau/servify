package auth

import "strings"

// HasPermission returns true if required is satisfied by granted permissions.
func HasPermission(granted []string, required string) bool {
	required = strings.TrimSpace(required)
	if required == "" {
		return true
	}
	for _, p := range granted {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if p == "*" || p == required {
			return true
		}
		if strings.HasSuffix(p, ".*") {
			prefix := strings.TrimSuffix(p, ".*")
			if prefix != "" && (required == prefix || strings.HasPrefix(required, prefix+".")) {
				return true
			}
		}
	}
	return false
}

func ResourcePermission(resource, method string) string {
	resource = strings.TrimSpace(resource)
	perm := resource + ".write"
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "GET", "HEAD", "OPTIONS":
		perm = resource + ".read"
	}
	return perm
}
