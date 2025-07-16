package apihandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/leonlonsdale/chirpy/internal/server/serverutil"
)

func RegisterValidateChirpsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/api/validate_chirp", validateChirpsHandler)
}

func validateChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		serverutil.RespondWithError(w, http.StatusInternalServerError, "Failed to decode request body", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		serverutil.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleaned := getCleanedChirp(params.Body, badWords)

	serverutil.RespondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleaned})
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
