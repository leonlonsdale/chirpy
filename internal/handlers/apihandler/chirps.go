package apihandler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterPostChirp(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /api/chirps", postChirpHandler(cfg))
}

func RegisterGetChirps(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /api/chirps", getAllChirpsHandler(cfg))
}

func RegisterGetChirpByID(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirpByIDHandler(cfg))
}

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func postChirpHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		params := parameters{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "could not decode json", err)
			return
		}

		if params.Body == "" {
			util.RespondWithError(w, http.StatusBadRequest, "a chirp must have a body", nil)
			return
		}

		if params.UserID == uuid.Nil {
			util.RespondWithError(w, http.StatusBadRequest, "a chirp must have an associated user id", nil)
			return
		}

		const maxChirpLength = 140
		if len(params.Body) > maxChirpLength {
			util.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}

		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		cleaned := getCleanedChirp(params.Body, badWords)

		data, err := cfg.DBQueries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: params.UserID,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating chirp", err)
			return
		}

		newChirp := chirp{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Body:      data.Body,
			UserID:    data.UserID,
		}

		util.RespondWithJSON(w, http.StatusCreated, newChirp)

	}
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

func getCleanedChirp(body string, badWords []string) string {

	words := strings.Split(body, " ")

	for i, w := range words {
		for _, bad := range badWords {
			if strings.EqualFold(w, bad) {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
