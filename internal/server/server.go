package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

func NewServer(port string, cfg *config.ApiConfig) *http.Server {

	mux := http.NewServeMux()

	registerHandlers(mux, cfg)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return s
}
