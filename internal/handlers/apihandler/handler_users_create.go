package apihandler

import (
	"encoding/json"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func CreateUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			handlers.User
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid json in request body", err)
			return
		}

		if params.Email == "" || params.Password == "" {
			util.RespondWithError(w, http.StatusBadRequest, "email and password are required", nil)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error hashing password", err)
			return
		}

		user, err := cfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
			return
		}

		newUser := response{
			User: handlers.User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt.Time,
				UpdatedAt:   user.UpdatedAt.Time,
				Email:       user.Email,
				IsChipryRed: user.IsChirpyRed,
			},
		}

		util.RespondWithJSON(w, http.StatusCreated, newUser)

	}
}
