package server

import (
	"fmt"
	"net/http"
)

func registerResetHandler(mux *http.ServeMux, cfg *apiConfig) {
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)

	_, _ = fmt.Fprintf(w, "Hit counter reset!")
}
