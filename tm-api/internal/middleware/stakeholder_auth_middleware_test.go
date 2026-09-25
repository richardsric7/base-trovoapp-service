package middleware

import "testing"

func TestStakeholderTokenFromHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
		ok     bool
	}{
		{name: "bearer", header: "Bearer token-1", want: "token-1", ok: true},
		{name: "lower bearer", header: "bearer token-2", want: "token-2", ok: true},
		{name: "raw", header: "token-3", want: "token-3", ok: true},
		{name: "empty", header: "", ok: false},
		{name: "malformed", header: "Bearer token extra", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stakeholderTokenFromHeader(tt.header)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("token = %q, want %q", got, tt.want)
			}
		})
	}
}
