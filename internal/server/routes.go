package server

import "net/http"

func registerHandlers(mux *http.ServeMux, cfg *ApiConfig) {
	registerHealthzHandler(mux)
	registerMetricsHandler(mux, cfg)
	registerResetHandler(mux, cfg)
	registerValidateChirpsHandler(mux)
}
