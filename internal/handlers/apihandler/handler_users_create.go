package apihandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterCreateUserHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /api/users", createUserHandler(cfg))
}

func createUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type registerResponseData struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid json in request body", err)
			return
		}

		if params.Email == "" {
			util.RespondWithError(w, http.StatusBadRequest, "email is required", nil)
			return
		}

		if params.Password == "" {
			util.RespondWithError(w, http.StatusBadRequest, "password is required", nil)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error hashing password", err)
			return
		}

		data, err := cfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		newUser := registerResponseData{
			ID:        data.ID.String(),
			CreatedAt: data.CreatedAt.Time,
			UpdatedAt: data.UpdatedAt.Time,
			Email:     data.Email,
		}

		util.RespondWithJSON(w, http.StatusCreated, newUser)

	}
}
