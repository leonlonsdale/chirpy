package server

import (
	"fmt"
	"net/http"
)

func registerHealthzHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", healthzHandler)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, http.StatusText(http.StatusOK))
}
