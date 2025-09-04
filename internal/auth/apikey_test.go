// apikey_test.go
package auth // <-- ändra till ditt paketnamn

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		auth        string
		wantKey     string
		wantIsErr   error  // använd för sentinel-felet
		wantErrText string // använd för “malformed authorization header”
	}{
		{
			name:      "missing header returns ErrNoAuthHeaderIncluded",
			auth:      "",
			wantIsErr: ErrNoAuthHeaderIncluded,
		},
		{
			name:    "happy path returns key",
			auth:    "ApiKey abc123",
			wantKey: "abc123",
		},
		{
			name:        "wrong scheme -> malformed",
			auth:        "Bearer abc123",
			wantErrText: "malformed authorization header",
		},
		{
			name:        "no value after scheme -> malformed",
			auth:        "ApiKey",
			wantErrText: "malformed authorization header",
		},
		{
			name:        "lowercase scheme -> malformed",
			auth:        "apikey abc123",
			wantErrText: "malformed authorization header",
		},
		{
			name:    "extra spaces after scheme (current behavior: empty key, no error)",
			auth:    "ApiKey   abc123",
			wantKey: "",
		},
		{
			name:    "multiple tokens: only first value returned",
			auth:    "ApiKey abc def",
			wantKey: "abc",
		},
		{
			name:    "trailing space after scheme (current behavior: empty key, no error)",
			auth:    "ApiKey ",
			wantKey: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := http.Header{}
			if tt.auth != "" {
				h.Set("Authorization", tt.auth)
			}

			got, err := GetAPIKey(h)

			// Kontroll av förväntade fel
			if tt.wantIsErr != nil {
				if !errors.Is(err, tt.wantIsErr) {
					t.Fatalf("expected error %v, got %v", tt.wantIsErr, err)
				}
				return
			}
			if tt.wantErrText != "" {
				if err == nil || err.Error() != tt.wantErrText {
					t.Fatalf("expected error %q, got %v", tt.wantErrText, err)
				}
				return
			}

			// Ingen error förväntad
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantKey {
				t.Fatalf("expected key %q, got %q", tt.wantKey, got)
			}
		})
	}
}
