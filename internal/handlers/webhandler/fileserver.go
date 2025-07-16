package webhandler

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func RegisterFileServerHandler(mux *http.ServeMux, cfg *serverconfig.ApiConfig) {
	fs := http.FileServer(http.Dir("./web/"))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", fs)))
}
