package webhandler

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

func FileServerHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))))
}
