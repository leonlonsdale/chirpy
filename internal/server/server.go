package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func NewServer(port string, cfg *serverconfig.ApiConfig) *http.Server {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/"))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(fileServer)))

	registerHandlers(mux, cfg)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return s
}
