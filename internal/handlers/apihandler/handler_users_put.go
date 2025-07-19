package apihandler

import (
	"encoding/json"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func UpdateUserHandler(db database.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type response struct {
			handlers.User
		}

		tokenStr, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
			return
		}

		userID, err := auth.ValidateJWT(tokenStr, secret)
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

		user, err := db.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "failed to update user", err)
			return
		}

		resp := response{
			User: handlers.User{
				ID:          user.ID,
				Email:       user.Email,
				CreatedAt:   user.CreatedAt.Time,
				UpdatedAt:   user.UpdatedAt.Time,
				IsChipryRed: user.IsChirpyRed,
			},
		}

		util.RespondWithJSON(w, http.StatusOK, resp)
	}
}
