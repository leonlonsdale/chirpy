package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/types"
)

type Chirps interface {
	Create(context.Context, types.NewChirp) (types.Chirp, error)
	GetAll(context.Context) ([]types.Chirp, error)
	GetById(context.Context, uuid.UUID) (types.Chirp, error)
	Delete(context.Context, types.DeleteChirp) error
}

type Users interface {
	Create(context.Context, types.NewUser) (types.User, error)
	GetByEmail(context.Context, string) (types.User, error)
	Update(context.Context, types.UpdateUser) (types.UpdatedUser, error)
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
		Chirps: &ChirpsStore{db},
		Users:  &UsersStore{db},
	}
}
