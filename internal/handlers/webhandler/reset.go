package webhandler

import (
	"fmt"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

func RegisterResetHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /admin/reset", ResetHandler(cfg))
}

func ResetHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.FileserverHits.Store(0)

		_, _ = fmt.Fprintf(w, "Hit counter reset!")
	}
}
