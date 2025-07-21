package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func (a *Auth) MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	n, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
	}
	if n != len(key) {
		return "", fmt.Errorf("insufficient random bytes read: expected %d, got %d", len(key), n)
	}
	return hex.EncodeToString(key), nil
}
