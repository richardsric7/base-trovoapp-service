package utils

import (
	"testing"
	"unicode"
)

func TestGenerateOTP_Table(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{
			name:        "valid_4_digits",
			length:      4,
			expectError: false,
		},
		{
			name:        "valid_6_digits",
			length:      6,
			expectError: false,
		},
		{
			name:        "valid_8_digits",
			length:      8,
			expectError: false,
		},
		{
			name:        "invalid_zero_length",
			length:      0,
			expectError: true,
		},
		{
			name:        "invalid_negative_length",
			length:      -3,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			otp, err := GenerateOTP(tc.length)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error for length %d, got nil", tc.length)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for length %d: %v", tc.length, err)
			}

			if len(otp) != tc.length {
				t.Errorf("expected length %d, got %d", tc.length, len(otp))
			}

			for _, ch := range otp {
				if !unicode.IsDigit(ch) {
					t.Errorf("expected only digits, got non-digit character %q", ch)
				}
			}
		})
	}
}
