package server

import (
	"fmt"
	"net/http"
)

func registerMetricsHandler(mux *http.ServeMux, cfg *apiConfig) {
	mux.HandleFunc("/metrics", cfg.metricsHandler)
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}
