package server

import "net/http"

type Server struct {
	mux    *http.ServeMux
	server *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/"))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	registerHandlers(mux)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &Server{
		mux:    mux,
		server: s,
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
