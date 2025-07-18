package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

type Server struct {
	Server *http.Server
	Mux    *http.ServeMux
}

func NewServer(port string, cfg *config.ApiConfig) *Server {

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	s := &Server{
		Server: server,
		Mux:    mux,
	}

	return s
}
