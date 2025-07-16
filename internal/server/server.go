package server

import (
	"net/http"
)

func NewServer(port string, cfg *ApiConfig) *http.Server {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/"))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServer)))

	registerHandlers(mux, cfg)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return s
}
