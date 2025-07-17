package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		wantToken   string
		wantErr     error
	}{
		{
			name:        "valid bearer token",
			headerValue: "Bearer abc123",
			wantToken:   "abc123",
			wantErr:     nil,
		},
		{
			name:        "no authorization header",
			headerValue: "",
			wantToken:   "",
			wantErr:     ErrorNoBearerToken,
		},
		{
			name:        "wrong scheme",
			headerValue: "Token abc123",
			wantToken:   "",
			wantErr:     ErrorNoBearerToken,
		},
		{
			name:        "missing token part",
			headerValue: "Bearer",
			wantToken:   "",
			wantErr:     ErrorNoBearerToken,
		},
		{
			name:        "extra spacing",
			headerValue: "Bearer     abc123",
			wantToken:   "abc123",
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerValue != "" {
				headers.Set("Authorization", tt.headerValue)
			}

			gotToken, gotErr := GetBearerToken(headers)
			if gotToken != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, gotToken)
			}
			if gotErr != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
		})
	}
}
