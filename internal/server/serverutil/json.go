package serverutil

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	resp := errorResponse{Error: msg}
	RespondWithJSON(w, code, resp)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
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
