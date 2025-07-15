package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func registerValidateChirpsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/api/validate_chirp", validateChirpsHandler)
}

func validateChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(params.Body, " ")

	for i, w := range words {
		for _, bad := range profaneWords {
			if strings.EqualFold(w, bad) {
				words[i] = "****"
				break
			}
		}
	}

	cleaned := strings.Join(words, " ")

	respBody := cleanedResponse{
		CleanedBody: cleaned,
	}

	respondWithJSON(w, http.StatusOK, respBody)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	resp := errorResponse{Error: msg}
	respondWithJSON(w, code, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}
