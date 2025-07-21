package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/storage"
	"github.com/leonlonsdale/chirpy/internal/types"
	"github.com/leonlonsdale/chirpy/internal/util"
)

type UserHandlers struct {
	cfg   *config.Config
	store *storage.Storage
}

func (h *UserHandlers) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		type response struct {
			User types.User
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "error decoding body", err)
			return
		}
		if params.Email == "" || params.Password == "" {
			util.RespondWithError(w, http.StatusBadRequest, "email and password are required", nil)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error hashing password", err)
			return
		}
		var newUser = types.NewUser{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		}

		user, err := h.store.Users.Create(ctx, newUser)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		resp := response{
			User: types.User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
				Email:       user.Email,
				IsChirpyRed: user.IsChirpyRed,
			},
		}

		util.RespondWithJSON(w, http.StatusCreated, resp)

	}
}

func (h *UserHandlers) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		type response struct {
			types.User
		}

		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
			return
		}

		userID, err := auth.ValidateJWT(tokenStr, h.cfg.Secret)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		type updateUserRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req updateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid request body", err)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error hashing password", err)
			return
		}
		var updates = types.UpdateUser{

			Email:          req.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		}

		user, err := h.store.Users.Update(ctx, updates)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "failed to update user", err)
			return
		}

		resp := response{
			User: types.User{
				ID:          user.ID,
				Email:       user.Email,
				CreatedAt:   user.CreatedAt.Time,
				UpdatedAt:   user.UpdatedAt.Time,
				IsChirpyRed: user.IsChirpyRed,
			},
		}

		util.RespondWithJSON(w, http.StatusOK, resp)
	}
}
