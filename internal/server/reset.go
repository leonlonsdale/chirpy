package server

import (
	"fmt"
	"net/http"
)

func registerResetHandler(mux *http.ServeMux, cfg *ApiConfig) {
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
}

func (cfg *ApiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)

	_, _ = fmt.Fprintf(w, "Hit counter reset!")
}
