package apihandler

import (
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterRefreshHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /api/refresh", refreshHandler(cfg))
}

func RegisterRevokeHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /api/revoke", revokeHandler(cfg))
}

func refreshHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		user, err := cfg.DBQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "couldn't find user for refresh token, the token may have expired", err)
			return
		}

		newAccessToken, err := auth.MakeJWT(user.ID, cfg.Secret, time.Hour)
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

func revokeHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		if err := cfg.DBQueries.RevokeRefreshToken(r.Context(), token); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem updating the refresh token record", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
