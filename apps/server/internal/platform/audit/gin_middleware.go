package audit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	beforeKey       = "audit.before"
	afterKey        = "audit.after"
	actionKey       = "audit.action"
	resourceTypeKey = "audit.resource_type"
	resourceIDKey   = "audit.resource_id"
)

func Middleware(recorder Recorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		if recorder == nil || !shouldAuditMethod(c.Request.Method) {
			c.Next()
			return
		}

		requestJSON := captureRequestBody(c)
		c.Next()

		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}

		entry := Entry{
			ActorUserID:   actorUserID(c),
			PrincipalKind: stringValue(c, "principal_kind"),
			Action:        contextOrDefault(c, actionKey, inferAction(c.FullPath(), c.Request.Method)),
			ResourceType:  contextOrDefault(c, resourceTypeKey, inferResourceType(c.FullPath())),
			ResourceID:    contextOrDefault(c, resourceIDKey, inferResourceID(c)),
			Route:         c.FullPath(),
			Method:        strings.ToUpper(strings.TrimSpace(c.Request.Method)),
			StatusCode:    c.Writer.Status(),
			Success:       c.Writer.Status() < http.StatusBadRequest,
			RequestID:     c.GetHeader("X-Request-ID"),
			ClientIP:      c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			TenantID:      stringValue(c, "tenant_id"),
			WorkspaceID:   stringValue(c, "workspace_id"),
			RequestJSON:   requestJSON,
			BeforeJSON:    contextJSON(c, beforeKey),
			AfterJSON:     contextJSON(c, afterKey),
		}

		if err := recorder.Record(c.Request.Context(), entry); err != nil {
			logrus.WithError(err).Warn("failed to persist audit log")
		}
	}
}

func SetBefore(c *gin.Context, value interface{}) {
	if c != nil {
		c.Set(beforeKey, value)
	}
}

func SetAfter(c *gin.Context, value interface{}) {
	if c != nil {
		c.Set(afterKey, value)
	}
}

func SetAction(c *gin.Context, value string) {
	if c != nil {
		c.Set(actionKey, strings.TrimSpace(value))
	}
}

func SetResourceType(c *gin.Context, value string) {
	if c != nil {
		c.Set(resourceTypeKey, strings.TrimSpace(value))
	}
}

func SetResourceID(c *gin.Context, value string) {
	if c != nil {
		c.Set(resourceIDKey, strings.TrimSpace(value))
	}
}

func shouldAuditMethod(method string) bool {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func captureRequestBody(c *gin.Context) string {
	if c == nil || c.Request == nil || c.Request.Body == nil {
		return ""
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return redactJSONText(strings.TrimSpace(string(body)))
}

func actorUserID(c *gin.Context) *uint {
	if c == nil {
		return nil
	}
	v, ok := c.Get("user_id")
	if !ok {
		return nil
	}
	id, ok := v.(uint)
	if !ok {
		return nil
	}
	return &id
}

func stringValue(c *gin.Context, key string) string {
	if c == nil {
		return ""
	}
	v, ok := c.Get(key)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func contextOrDefault(c *gin.Context, key, fallback string) string {
	if value := stringValue(c, key); value != "" {
		return value
	}
	return fallback
}

func contextJSON(c *gin.Context, key string) string {
	if c == nil {
		return ""
	}
	v, ok := c.Get(key)
	if !ok || v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return redactJSONText(string(data))
}

func redactJSONText(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return raw
	}

	redacted, changed := redactValue(payload)
	if !changed {
		return raw
	}

	data, err := json.Marshal(redacted)
	if err != nil {
		return raw
	}
	return string(data)
}

func redactValue(value interface{}) (interface{}, bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		changed := false
		for key, item := range typed {
			if shouldRedactKey(key) {
				typed[key] = "[REDACTED]"
				changed = true
				continue
			}
			next, nestedChanged := redactValue(item)
			if nestedChanged {
				typed[key] = next
				changed = true
			}
		}
		return typed, changed
	case []interface{}:
		changed := false
		for i, item := range typed {
			next, nestedChanged := redactValue(item)
			if nestedChanged {
				typed[i] = next
				changed = true
			}
		}
		return typed, changed
	default:
		return value, false
	}
}

func shouldRedactKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")

	switch normalized {
	case "password", "passwd", "secret", "clientsecret", "apikey", "token", "accesstoken", "refreshtoken", "authorization":
		return true
	default:
		return strings.HasSuffix(normalized, "secret") || strings.HasSuffix(normalized, "token") || strings.HasSuffix(normalized, "apikey")
	}
}

func inferResourceType(route string) string {
	parts := routeParts(route)
	if len(parts) == 0 {
		return "unknown"
	}
	if parts[0] == "api" {
		parts = parts[1:]
	}
	if len(parts) > 0 && parts[0] == "v1" {
		parts = parts[1:]
	}
	if len(parts) == 0 {
		return "unknown"
	}
	switch parts[0] {
	case "omni":
		if len(parts) > 1 {
			return sanitizeSegment(parts[1])
		}
	case "apps":
		if len(parts) > 1 {
			return "app_" + sanitizeSegment(parts[1])
		}
	}
	return sanitizeSegment(parts[0])
}

func inferAction(route, method string) string {
	parts := routeParts(route)
	if len(parts) == 0 {
		return strings.ToLower(method)
	}
	resource := inferResourceType(route)
	last := ""
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.HasPrefix(parts[i], ":") {
			continue
		}
		last = sanitizeSegment(parts[i])
		break
	}
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodPost:
		if last == "" || last == resource {
			return resource + ".create"
		}
		return resource + "." + last
	case http.MethodPut, http.MethodPatch:
		if last == "" || last == resource {
			return resource + ".update"
		}
		return resource + "." + last
	case http.MethodDelete:
		return resource + ".delete"
	default:
		return resource + "." + strings.ToLower(method)
	}
}

func inferResourceID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	for _, key := range []string{"id", "ticket_id", "session_id", "recordingID", "recording_id", "protocol"} {
		if v := strings.TrimSpace(c.Param(key)); v != "" {
			return v
		}
	}
	return ""
}

func routeParts(route string) []string {
	trimmed := strings.Trim(strings.TrimSpace(route), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func sanitizeSegment(segment string) string {
	segment = strings.TrimSpace(segment)
	segment = strings.TrimPrefix(segment, ":")
	segment = strings.ReplaceAll(segment, "-", "_")
	return segment
}
