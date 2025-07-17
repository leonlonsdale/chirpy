package apihandler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterCreateChirpHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.Handle("POST /api/chirps", auth.AuthMiddleware(cfg, createChirpHandler(cfg)))
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func createChirpHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type response struct {
			Chirp
		}

		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok || userID == uuid.Nil {
			util.RespondWithError(w, http.StatusUnauthorized, "user not authenticated", nil)
			return
		}

		var params struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "could not decode json", err)
			return
		}

		if params.Body == "" {
			util.RespondWithError(w, http.StatusBadRequest, "a chirp must have a body", nil)
			return
		}

		if userID == uuid.Nil {
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
			UserID: userID,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating chirp", err)
			return
		}

		chirp := response{
			Chirp: Chirp{
				ID:        data.ID,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.UpdatedAt,
				Body:      data.Body,
				UserID:    data.UserID,
			},
		}

		util.RespondWithJSON(w, http.StatusCreated, chirp)

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
