package api

import "testing"

func TestIsValidMonth(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		valid bool
	}{
		{name: "valid", in: "2026-05", valid: true},
		{name: "invalid format", in: "2026/05", valid: false},
		{name: "invalid month", in: "2026-13", valid: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidMonth(tc.in); got != tc.valid {
				t.Fatalf("isValidMonth(%q) = %v, want %v", tc.in, got, tc.valid)
			}
		})
	}
}
