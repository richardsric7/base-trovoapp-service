package utils

import "testing"

func TestCallbackURLs(t *testing.T) {
	for _, base := range []string{"https://tm.example.com/api/v1/callbacks/login", "https://tm.example.com/api/v1/callbacks/login/"} {
		t.Setenv("LOGIN_CALLBACK_URL", base)
		if got, want := LoginCallbackURL("trovo"), "https://tm.example.com/api/v1/callbacks/login/trovo"; got != want {
			t.Errorf("LoginCallbackURL(%q) = %q, want %q", base, got, want)
		}
		if got, want := AuthorizationCallbackURL("trovo"), "https://tm.example.com/api/v1/callbacks/auth/trovo"; got != want {
			t.Errorf("AuthorizationCallbackURL(%q) = %q, want %q", base, got, want)
		}
	}
}

func TestFrontendBaseURLHasNoFallback(t *testing.T) {
	t.Setenv("TROVO_MANAGER_BASE_URL", "")
	if got := GetTrovomanagerFrontendBaseUrl(); got != "" {
		t.Errorf("want no built-in fallback, got %q", got)
	}
}
