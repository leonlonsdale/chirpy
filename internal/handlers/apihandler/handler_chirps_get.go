package apihandler

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterGetAllChirpsHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /api/chirps", getAllChirpsHandler(cfg))
}

func RegisterGetChirpByIDHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirpByIDHandler(cfg))
}

func getAllChirpsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpsData, err := cfg.DBQueries.GetAllChirps(r.Context())
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
		}

		chirps := make([]chirp, 0, len(chirpsData))
		for _, c := range chirpsData {
			chirps = append(chirps, chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}

		util.RespondWithJSON(w, http.StatusOK, chirps)

	}
}

func getChirpByIDHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid chirp id format", err)
		}
		chirpData, err := cfg.DBQueries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			if err == sql.ErrNoRows {
				util.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
				return
			}

			util.RespondWithError(w, http.StatusInternalServerError, "error fetching chirp", err)
			return

		}

		foundChirp := chirp{
			ID:        chirpData.ID,
			CreatedAt: chirpData.CreatedAt,
			UpdatedAt: chirpData.UpdatedAt,
			Body:      chirpData.Body,
			UserID:    chirpData.UserID,
		}

		util.RespondWithJSON(w, http.StatusOK, foundChirp)

	}
}
