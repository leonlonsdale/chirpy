package webhandler

import (
	"fmt"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func RegisterResetHandler(mux *http.ServeMux, cfg *serverconfig.ApiConfig) {
	mux.HandleFunc("POST /admin/reset", ResetHandler(cfg))
}

func ResetHandler(cfg *serverconfig.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.FileserverHits.Store(0)

		_, _ = fmt.Fprintf(w, "Hit counter reset!")
	}
}
