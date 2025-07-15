package server

import "net/http"

func registerHandlers(mux *http.ServeMux, cfg *apiConfig) {
	registerHealthzHandler(mux)
	registerMetricsHandler(mux, cfg)
	registerResetHandler(mux, cfg)
}
