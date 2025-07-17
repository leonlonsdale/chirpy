package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTCreateAndValidate_Success(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if returnedID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, returnedID)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -1*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestJWT_InvalidSecret(t *testing.T) {
	secret := "correct-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error for invalid secret, got nil")
	}
}
