package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterUpdateUserHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("PUT /api/users", updateUserHandler(cfg))
}

func updateUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse access token
		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
			return
		}

		// Validate token and extract user ID
		subject, err := auth.ValidateJWT(tokenStr, cfg.Secret)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		// Parse request body
		type updateUserRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req updateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
			return
		}

		// Hash new password
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error hashing password", err)
			return
		}

		// Update user in DB
		user, err := cfg.DBQueries.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
			ID:             subject,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "failed to update user", err)
			return
		}
		type userResponse struct {
			ID        uuid.UUID `json:"id"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}

		resp := userResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		}

		// Respond with updated user (without password)
		util.RespondWithJSON(w, http.StatusOK, resp)
	}
}
