package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/storage"
	"github.com/leonlonsdale/chirpy/internal/types"
	"github.com/leonlonsdale/chirpy/internal/util"
)

type ChirpHandlers struct {
	Store *storage.Storage
	Cfg   *config.Config
}

func (h *ChirpHandlers) CreateChirp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			types.Chirp
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

		data, err := h.Store.Chirps.Create(r.Context(), types.NewChirp{
			Body:   cleaned,
			UserID: userID,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error creating chirp", err)
			return
		}

		chirp := response{
			Chirp: types.Chirp{
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

func (h *ChirpHandlers) GetAllChirps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authorID, filterByAuthor, err := parseAuthorID(r)
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid author_id", err)
			return
		}

		sortOrder := r.URL.Query().Get("sort")
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}

		chirpsData, err := h.Store.Chirps.GetAll(ctx)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
			return
		}

		filteredChirps := filterChirpsByAuthor(chirpsData, authorID, filterByAuthor)
		sortChirps(filteredChirps, sortOrder)
		util.RespondWithJSON(w, http.StatusOK, filteredChirps)
	}
}

func (h *ChirpHandlers) GetChirpById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type response struct {
			types.Chirp
		}
		ctx := r.Context()
		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			http.Error(w, "Invalid Chirp ID", http.StatusBadRequest)
			return
		}

		if chirpID == uuid.Nil {
			http.Error(w, "Chirp ID is required", http.StatusBadRequest)
			return
		}

		chirp, err := h.Store.Chirps.GetById(ctx, chirpID)
		if err != nil {
			if err == sql.ErrNoRows {
				util.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
				return
			}
			http.Error(w, "Failed to retrieve chirp", http.StatusInternalServerError)
			return
		}
		foundChirp := response{
			Chirp: types.Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			},
		}

		util.RespondWithJSON(w, http.StatusOK, foundChirp)

	}
}

func (h *ChirpHandlers) DeleteChirpById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid jwt", nil)
			return
		}

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid chirp id format", err)
			return
		}

		chirp, err := h.Store.Chirps.GetById(r.Context(), chirpID)
		if err != nil {
			util.RespondWithError(w, http.StatusNotFound, "could not find chirp to delete", err)
			return
		}

		if userID != chirp.UserID {
			util.RespondWithError(w, http.StatusForbidden, "not author of chirp", nil)
			return
		}

		if err := h.Store.Chirps.Delete(r.Context(), types.DeleteChirp{
			ID:     chirpID,
			UserID: userID,
		}); err != nil {
			util.RespondWithError(w, http.StatusNotFound, "error deleting chirp", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func parseAuthorID(r *http.Request) (authorID uuid.UUID, filterByAuthor bool, error error) {
	authorIDStr := r.URL.Query().Get("author_id")
	if authorIDStr == "" {
		return uuid.Nil, false, nil
	}
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return uuid.Nil, false, err
	}
	return authorID, true, nil
}

func filterChirpsByAuthor(chirps []types.Chirp, authorID uuid.UUID, filterByAuthor bool) []types.Chirp {
	filtered := make([]types.Chirp, 0, len(chirps))
	for _, c := range chirps {
		if filterByAuthor && c.UserID != authorID {
			continue
		}
		filtered = append(filtered, types.Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	return filtered
}

func sortChirps(chirps []types.Chirp, sortOrder string) {
	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "asc" {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})
}
