package middleware

import "os"

// AuthChecksDisabled reports whether request signature and API-key checks
// are switched off (ENABLE_AUTH_MIDDLEWARE=0). With them off anyone can act
// as any user, so the switch only works outside release mode: main refuses
// to start with ENABLE_AUTH_MIDDLEWARE=0 and GIN_MODE=release, and this
// ignores the switch there too, in case GIN_MODE changes at run time.
func AuthChecksDisabled() bool {
	return os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "0" && os.Getenv("GIN_MODE") != "release"
}

// MaskSecret shortens an API key or other secret for logs: enough to tell
// keys apart, never enough to use one.
func MaskSecret(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "…" + s[len(s)-2:]
}
