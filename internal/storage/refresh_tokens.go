package storage

import (
	"context"
	"database/sql"

	"github.com/leonlonsdale/chirpy/internal/types"
)

type RefreshTokenStore struct {
	db *sql.DB
}

func (rt *RefreshTokenStore) Create(ctx context.Context, arg types.CreateRefreshToken) error {
	query := `
	INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
	    VALUES ($1, NOW(), NOW(), $2, $3)
	RETURNING
	    token, created_at, updated_at, user_id, expires_at, revoked_at
	`
	_, err := rt.db.ExecContext(ctx, query, arg.Token, arg.UserID, arg.ExpiresAt)
	return err

}

func (rt *RefreshTokenStore) Get(ctx context.Context, token string) (types.RefreshToken, error) {
	query := `
		UPDATE
		    refresh_tokens
		SET
		    revoked_at = NOW(),
		    updated_at = NOW()
		WHERE
		    token = $1
		RETURNING
		    token, created_at, updated_at, user_id, expires_at, revoked_at
	`
	row := rt.db.QueryRowContext(ctx, query, token)
	var r types.RefreshToken
	err := row.Scan(
		&r.Token,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.UserID,
		&r.ExpiresAt,
		&r.RevokedAt,
	)
	return r, err

}

func (rt *RefreshTokenStore) GetUserFromToken(ctx context.Context, token string) (types.User, error) {
	query := `
		SELECT
		    users.id, users.created_at, users.updated_at, users.email, users.hashed_password, users.is_chirpy_red
		FROM
		    users
		    JOIN refresh_tokens ON users.id = refresh_tokens.user_id
		WHERE
		    refresh_tokens.token = $1
		    AND revoked_at IS NULL
		    AND expires_at > NOW()
	`
	row := rt.db.QueryRowContext(ctx, query, token)
	var r types.User
	err := row.Scan(
		&r.ID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Email,
		&r.HashedPassword,
		&r.IsChirpyRed,
	)
	return r, err
}

func (rt *RefreshTokenStore) Revoke(ctx context.Context, token string) error {
	query := `
		UPDATE
		    refresh_tokens
		SET
		    revoked_at = NOW(),
		    updated_at = NOW()
		WHERE
		    token = $1
	`
	_, err := rt.db.ExecContext(ctx, query, token)
	return err
}
