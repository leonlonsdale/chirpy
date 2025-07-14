package server

import "net/http"

func registerHandlers(mux *http.ServeMux) {
	registerHealthzHandler(mux)
}
