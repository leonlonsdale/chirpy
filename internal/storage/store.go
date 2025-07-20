package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Chirps interface {
	Create(context.Context, NewChirp) (Chirp, error)
	GetAll(context.Context) ([]Chirp, error)
	GetById(context.Context, uuid.UUID) (Chirp, error)
	Delete(context.Context, DeleteChirp) error
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
		Chirps: &ChirpsStore{db},
		Users:  &UsersStore{db},
	}
}
