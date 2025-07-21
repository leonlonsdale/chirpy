package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	HashPassword(string) (string, error)
	CheckPasswordHash(string, string) error
}

type Auth struct{}

func NewAuthService() Auth {
	return Auth{}
}

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrorJWTInvalidSigningMethod = errors.New("invalid signing method")
var ErrorJWTInvalidToken = errors.New("invalid jwt token")
var ErrorJWTInvalidClaimsType = errors.New("invalid claims type")
var ErrorJWTTokenHasExpired = errors.New("jwt token has expired")

var a2params = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Argon 2id

func (a *Auth) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, a2params)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return hash, nil
}

func (a *Auth) CheckPasswordHash(password, hash string) error {

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("Error: password hash compare: %v\n", err)
		return err
	}

	if !match {
		return ErrInvalidCredentials
	}

	return nil
}

// JWT

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func (a *Auth) MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

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

func (a *Auth) ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
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
