package utils

import (
	"os"
	"strings"
)

// GetTrovomanagerFrontendBaseUrl is tm-web's address, for links in emails.
// TROVO_MANAGER_BASE_URL is required at start (main.go).
func GetTrovomanagerFrontendBaseUrl() string {
	return os.Getenv("TROVO_MANAGER_BASE_URL")
}

// LoginCallbackURL is where app-backend reports an approved login:
// LOGIN_CALLBACK_URL (https://<tm-api>/api/v1/callbacks/login) plus the
// service name, matching the /api/v1/callbacks/login/:serviceName route.
func LoginCallbackURL(serviceName string) string {
	return strings.TrimRight(os.Getenv("LOGIN_CALLBACK_URL"), "/") + "/" + serviceName
}

// AuthorizationCallbackURL is where app-backend reports an approved
// authorization request: the /api/v1/callbacks/auth/:serviceName route,
// derived from LOGIN_CALLBACK_URL so no extra setting is needed.
func AuthorizationCallbackURL(serviceName string) string {
	base := strings.TrimRight(os.Getenv("LOGIN_CALLBACK_URL"), "/")
	if strings.HasSuffix(base, "/callbacks/login") {
		base = strings.TrimSuffix(base, "/callbacks/login") + "/callbacks/auth"
	}
	return base + "/" + serviceName
}
