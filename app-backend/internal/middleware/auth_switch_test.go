package middleware

import "testing"

func TestAuthChecksDisabledNeverInRelease(t *testing.T) {
	cases := []struct {
		enable, mode string
		want         bool
	}{
		{"", "", false},
		{"1", "", false},
		{"0", "", true},
		{"0", "debug", true},
		{"0", "release", false},
	}
	for _, c := range cases {
		t.Setenv("ENABLE_AUTH_MIDDLEWARE", c.enable)
		t.Setenv("GIN_MODE", c.mode)
		if got := AuthChecksDisabled(); got != c.want {
			t.Errorf("ENABLE_AUTH_MIDDLEWARE=%q GIN_MODE=%q: got %v, want %v", c.enable, c.mode, got, c.want)
		}
	}
}
