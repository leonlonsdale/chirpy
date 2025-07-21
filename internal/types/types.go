package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// user

type NewUser struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

type UpdateUser struct {
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	ID             uuid.UUID `json:"id"`
}

type UpdatedUser struct {
	ID          uuid.UUID    `json:"id"`
	Email       string       `json:"email"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	IsChirpyRed bool         `json:"is_chirpy_red"`
}

// chirps

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type NewChirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type DeleteChirp struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// refresh token

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

type CreateRefreshToken struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
