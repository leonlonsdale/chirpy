package webhandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func RegisterResetHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("POST /admin/reset", ResetHandler(cfg))
}

func ResetHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Platform != "dev" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			log.Println("User accessed reset")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.FileserverHits.Store(0)
		if err := cfg.DBQueries.ResetUsers(r.Context()); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem resetting users", err)
		}
		_, _ = fmt.Fprintf(w, "Hit counter reset!")
	}
}
