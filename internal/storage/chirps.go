package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/pkg/utils"
)

type ChirpsStore struct {
	db *sql.DB
}

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

func (cs *ChirpsStore) Create(ctx context.Context, data NewChirp) (Chirp, error) {
	query := `
		INSERT INTO chirps (id, created_at, updated_at, body, user_id)
		    VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
		RETURNING
		    id, created_at, updated_at, body, user_id
	`

	row := cs.db.QueryRowContext(ctx, query, data.Body, data.UserID)
	var c Chirp
	err := row.Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Body,
		&c.UserID,
	)

	return c, err
}

func (cs *ChirpsStore) Delete(ctx context.Context, data DeleteChirp) error {
	query := `
		DELETE FROM chirps
		WHERE id = $1
		    AND user_id = $2
	`

	_, err := cs.db.ExecContext(ctx, query)
	return err
}

func (cs *ChirpsStore) GetAll(ctx context.Context) ([]Chirp, error) {
	query := `
		SELECT
		    id, created_at, updated_at, body, user_id
		FROM
		    chirps
		ORDER BY
		    created_at
	`

	rows, err := cs.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer utils.SafeClose(rows)

	var items []Chirp
	for rows.Next() {
		var c Chirp
		if err := rows.Scan(
			&c.ID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Body,
			&c.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (cs *ChirpsStore) GetById(ctx context.Context, id uuid.UUID) (Chirp, error) {
	query := `
		SELECT
		    id, created_at, updated_at, body, user_id
		FROM
		    chirps
		WHERE
		    id = $1
	`
	row := cs.db.QueryRowContext(ctx, query, id)
	var c Chirp
	err := row.Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Body,
		&c.UserID,
	)
	return c, err
}
