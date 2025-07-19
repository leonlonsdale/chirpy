package apihandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
)

const validEvent = "user.upgraded"

func UpgradeToChirpyRedHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != cfg.PolkaKey {
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

		affectedUsers, err := cfg.DBQueries.UpgradeToChirpyRed(r.Context(), userID)
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
