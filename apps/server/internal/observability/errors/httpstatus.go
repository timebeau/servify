package errors

import "net/http"

// HTTPStatusFromError returns the suggested HTTP status code for any error.
// Returns 500 for unclassified errors.
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	appErr := Classify(err)
	if appErr != nil && appErr.HTTPStatus > 0 {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// UserMessageFromError returns a user-safe message for any error.
// Returns "Internal server error" for unclassified errors to avoid leaking details.
func UserMessageFromError(err error) string {
	if err == nil {
		return ""
	}
	appErr := Classify(err)
	if appErr != nil && appErr.Message != "" {
		// For system/config errors, return generic message
		if appErr.Severity == SeveritySystem || appErr.Severity == SeverityConfig {
			return "Internal server error"
		}
		return appErr.Message
	}
	return "Internal server error"
}
