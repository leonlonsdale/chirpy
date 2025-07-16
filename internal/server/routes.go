package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/handlers/apihandler"
	"github.com/leonlonsdale/chirpy/internal/handlers/webhandler"
	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func registerHandlers(mux *http.ServeMux, cfg *serverconfig.ApiConfig) {
	webhandler.RegisterMetricsHandler(mux, cfg)
	webhandler.RegisterResetHandler(mux, cfg)
	webhandler.RegisterFileServerHandler(mux, cfg)

	apihandler.RegisterHealthzHandler(mux)
	apihandler.RegisterValidateChirpsHandler(mux)
}
