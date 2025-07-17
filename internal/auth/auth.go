package auth

import (
	"errors"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

var params = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) error {

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
