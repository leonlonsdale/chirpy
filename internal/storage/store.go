package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/types"
)

type Chirps interface {
	Create(context.Context, types.NewChirp) (database.Chirp, error)
	Read(context.Context) error
	Update(context.Context) error
	Delete(context.Context) error
}

type Users interface {
	Create(context.Context, NewUser) (User, error)
	GetByEmail(context.Context, string) (User, error)
	Update(context.Context, UpdateUser) (UpdatedUser, error)
	Delete(context.Context) error
	Upgrade(context.Context, uuid.UUID) (int64, error)
	Reset(context.Context) error
}

type Storage struct {
	Users  Users
	Chirps Chirps
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		// Chirps: &ChirpsStore{db},
		Users: &UsersStore{db},
	}
}
