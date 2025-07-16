package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterCreateUserHandler(mux *http.ServeMux, cfg *config.ApiConfig) {

	mux.HandleFunc("POST /api/users", createUserHandler(cfg))
}

func createUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Email string `json:"email"`
		}

		type resValue struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}

		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to decode request body", err)
			return
		}

		if params.Email == "" {
			util.RespondWithError(w, http.StatusBadRequest, "Email is required", nil)
			return
		}

		data, err := cfg.DBQueries.CreateUser(r.Context(), params.Email)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		}

		newUser := resValue{
			ID:        data.ID.String(),
			CreatedAt: data.CreatedAt.Time,
			UpdatedAt: data.UpdatedAt.Time,
			Email:     data.Email,
		}

		util.RespondWithJSON(w, http.StatusCreated, newUser)

	}
}
