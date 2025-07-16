package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func NewServer(port string, cfg *serverconfig.ApiConfig) *http.Server {

	mux := http.NewServeMux()

	registerHandlers(mux, cfg)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return s
}
