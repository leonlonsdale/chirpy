package apihandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/database"
)

const validEvent = "user.upgraded"

func UpgradeToChirpyRedHandler(db database.Queries, polkakey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != polkakey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if params.Event != validEvent {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		userID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		affectedUsers, err := db.UpgradeToChirpyRed(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if affectedUsers == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
