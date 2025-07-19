package apihandler

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func GetAllChirpsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorIDStr := r.URL.Query().Get("author_id")
		var authorID uuid.UUID
		var err error
		filterByAuthor := false

		if authorIDStr != "" {
			authorID, err = uuid.Parse(authorIDStr)
			if err != nil {
				util.RespondWithError(w, http.StatusBadRequest, "invalid author_id", err)
				return
			}
			filterByAuthor = true
		}

		sortOrder := r.URL.Query().Get("sort")
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}

		chirpsData, err := cfg.DBQueries.GetAllChirps(r.Context())
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
			return
		}

		filteredChirps := make([]handlers.Chirp, 0, len(chirpsData))
		for _, c := range chirpsData {
			if filterByAuthor && c.UserID != authorID {
				continue
			}
			filteredChirps = append(filteredChirps, handlers.Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}

		sort.Slice(filteredChirps, func(i, j int) bool {
			if sortOrder == "asc" {
				return filteredChirps[i].CreatedAt.Before(filteredChirps[j].CreatedAt)
			}
			return filteredChirps[i].CreatedAt.After(filteredChirps[j].CreatedAt)
		})

		util.RespondWithJSON(w, http.StatusOK, filteredChirps)
	}
}

func GetChirpByIDHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type response struct {
			handlers.Chirp
		}

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

		foundChirp := response{
			Chirp: handlers.Chirp{
				ID:        chirpData.ID,
				CreatedAt: chirpData.CreatedAt,
				UpdatedAt: chirpData.UpdatedAt,
				Body:      chirpData.Body,
				UserID:    chirpData.UserID,
			},
		}

		util.RespondWithJSON(w, http.StatusOK, foundChirp)

	}
}
