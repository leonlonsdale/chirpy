package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrorJWTInvalidSigningMethod = errors.New("invalid signing method")
var ErrorJWTInvalidToken = errors.New("invalid jwt token")
var ErrorJWTInvalidClaimsType = errors.New("invalid claims type")
var ErrorJWTTokenHasExpired = errors.New("jwt token has expired")

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	keyfunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, ErrorJWTInvalidSigningMethod
		}
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token string: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, ErrorJWTInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, ErrorJWTInvalidClaimsType
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now().UTC()) {
		return uuid.Nil, ErrorJWTTokenHasExpired
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user id not found in token validation: %w", err)
	}

	return uid, nil
}
