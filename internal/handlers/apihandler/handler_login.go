package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterLoginHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /api/login", loginHandler(cfg))
}

func loginHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type loginResponse struct {
			ID           string    `json:"id"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Email        string    `json:"email"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid json in request body", err)
			return
		}

		userData, err := cfg.DBQueries.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error looking up user", err)
			return
		}

		if err := auth.CheckPasswordHash(params.Password, userData.HashedPassword); err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
			return
		}

		token, err := auth.MakeJWT(userData.ID, cfg.Secret, time.Hour)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating jwt token", err)
			return
		}

		rToken, err := auth.MakeRefreshToken()
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
			return
		}

		err = cfg.DBQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     rToken,
			UserID:    userData.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error saving refresh token", err)
			return
		}

		respData := loginResponse{
			ID:           userData.ID.String(),
			CreatedAt:    userData.CreatedAt.Time,
			UpdatedAt:    userData.UpdatedAt.Time,
			Email:        userData.Email,
			Token:        token,
			RefreshToken: rToken,
		}

		util.RespondWithJSON(w, http.StatusOK, respData)

	}
}
