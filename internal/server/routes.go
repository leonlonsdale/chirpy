package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers/apihandler"
	"github.com/leonlonsdale/chirpy/internal/handlers/webhandler"
)

func registerHandlers(mux *http.ServeMux, cfg *config.ApiConfig) {
	webhandler.RegisterMetricsHandler(mux, cfg)
	webhandler.RegisterResetHandler(mux, cfg)
	webhandler.RegisterFileServerHandler(mux, cfg)

	apihandler.RegisterHealthzHandler(mux)
	apihandler.RegisterValidateChirpsHandler(mux)
}
