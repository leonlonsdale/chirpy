package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/storage"
	"github.com/leonlonsdale/chirpy/internal/types"
	"github.com/leonlonsdale/chirpy/internal/util"
)

type AuthHandlers struct {
	store *storage.Storage
	cfg   *config.Config
	auth  auth.Auth
}

func (ah *AuthHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			types.User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid json in request body", err)
			return
		}

		user, err := ah.store.Users.GetByEmail(r.Context(), params.Email)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error looking up user", err)
			return
		}

		if err := ah.auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
			return
		}

		accessToken, err := ah.auth.MakeJWT(user.ID, ah.cfg.Secret, time.Hour)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating jwt token", err)
			return
		}

		refreshToken, err := ah.auth.MakeRefreshToken()
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
			return
		}
		var tokenParams = types.CreateRefreshToken{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		}

		if err := ah.store.RefreshToken.Create(r.Context(), tokenParams); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error saving refresh token", err)
			return
		}

		responseBody := response{
			User: types.User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
				Email:       user.Email,
				IsChirpyRed: user.IsChirpyRed,
			},
			Token:        accessToken,
			RefreshToken: refreshToken,
		}

		util.RespondWithJSON(w, http.StatusOK, responseBody)

	}
}

func (ah *AuthHandlers) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		refreshToken, err := ah.auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		user, err := ah.store.RefreshToken.GetUserFromToken(r.Context(), refreshToken)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "couldn't find user for refresh token, the token may have expired", err)
			return
		}

		newAccessToken, err := ah.auth.MakeJWT(user.ID, ah.cfg.Secret, time.Hour)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating new access token", err)
			return
		}

		refreshResp := response{
			Token: newAccessToken,
		}

		util.RespondWithJSON(w, http.StatusOK, refreshResp)

	}
}

func (ah *AuthHandlers) Revoke() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ah.auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		if err := ah.store.RefreshToken.Revoke(r.Context(), token); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem updating the refresh token record", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
