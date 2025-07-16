package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/handlers/apihandler"
	"github.com/leonlonsdale/chirpy/internal/handlers/webhandler"
	"github.com/leonlonsdale/chirpy/internal/server/serverconfig"
)

func registerHandlers(mux *http.ServeMux, cfg *serverconfig.ApiConfig) {
	apihandler.RegisterHealthzHandler(mux)
	webhandler.RegisterMetricsHandler(mux, cfg)
	webhandler.RegisterResetHandler(mux, cfg)
	apihandler.RegisterValidateChirpsHandler(mux)
}
