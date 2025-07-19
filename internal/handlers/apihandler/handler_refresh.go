package apihandler

import (
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RefreshHandler(db database.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		user, err := db.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "couldn't find user for refresh token, the token may have expired", err)
			return
		}

		newAccessToken, err := auth.MakeJWT(user.ID, secret, time.Hour)
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

func RevokeHandler(db database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve bearer token", err)
			return
		}

		if err := db.RevokeRefreshToken(r.Context(), token); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem updating the refresh token record", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
