package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// user

type NewUser struct {
	Email          string
	HashedPassword string
}

type User struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsChirpyRed    bool
}

type UpdateUser struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

type UpdatedUser struct {
	ID          uuid.UUID
	Email       string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
	IsChirpyRed bool
}

// chirps

type Chirp struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	UserID    uuid.UUID
}

type NewChirp struct {
	Body   string
	UserID uuid.UUID
}

type DeleteChirp struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
