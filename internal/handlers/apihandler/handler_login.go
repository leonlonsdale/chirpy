package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func LoginHandler(db database.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			handlers.User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid json in request body", err)
			return
		}

		user, err := db.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error looking up user", err)
			return
		}

		if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
			return
		}

		accessToken, err := auth.MakeJWT(user.ID, secret, time.Hour)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating jwt token", err)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
			return
		}

		if err := db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		}); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error saving refresh token", err)
			return
		}

		responseBody := response{
			User: handlers.User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt.Time,
				UpdatedAt:   user.UpdatedAt.Time,
				Email:       user.Email,
				IsChipryRed: user.IsChirpyRed,
			},
			Token:        accessToken,
			RefreshToken: refreshToken,
		}

		util.RespondWithJSON(w, http.StatusOK, responseBody)

	}
}
