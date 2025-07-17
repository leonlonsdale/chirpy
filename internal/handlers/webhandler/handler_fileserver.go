package webhandler

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

func RegisterFileServerHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	fs := http.FileServer(http.Dir("./web/"))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", fs)))
}
