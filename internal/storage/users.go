package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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

type UsersStore struct {
	db *sql.DB
}

func (us *UsersStore) Create(ctx context.Context, data NewUser) (User, error) {
	query := `
		INSERT INTO users (id, created_at, updated_at, email, hashed_password)
    		VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
		RETURNING
    		id, created_at, updated_at, email, hashed_password, is_chirpy_red
	`
	row := us.db.QueryRowContext(ctx, query, data.Email, data.HashedPassword)

	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.HashedPassword,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsChirpyRed,
	)

	return u, err
}

func (us *UsersStore) GetByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT
		    id, created_at, updated_at, email, hashed_password, is_chirpy_red
		FROM
		    users
		WHERE
		    email = $1
	`

	row := us.db.QueryRowContext(ctx, query, email)
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.HashedPassword,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsChirpyRed,
	)

	return u, err
}

func (us *UsersStore) Update(ctx context.Context, data UpdateUser) (UpdatedUser, error) {
	query := `
		UPDATE
		    users
		SET
		    email = $1,
		    hashed_password = $2,
		    updated_at = NOW()
		WHERE
		    id = $3
		RETURNING
		    id,
		    email,
		    created_at,
		    updated_at,
		    is_chirpy_red
	`
	row := us.db.QueryRowContext(ctx, query, data.Email, data.HashedPassword, data.ID)
	var u UpdatedUser
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsChirpyRed,
	)

	return u, err
}

func (us *UsersStore) Delete(ctx context.Context) error {
	return nil
}

func (us *UsersStore) Upgrade(ctx context.Context, id uuid.UUID) (int64, error) {
	query := `
		UPDATE
		    users
		SET
		    is_chirpy_red = TRUE
		WHERE
		    id = $1
	`

	result, err := us.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func (us *UsersStore) Reset(ctx context.Context) error {
	query := `
		DELETE FROM users
	`

	_, err := us.db.ExecContext(ctx, query)
	return err
}
