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

	// ==========[ API ]
	apihandler.RegisterHealthzHandler(mux)

	// user
	apihandler.RegisterCreateUserHandler(mux, cfg)
	apihandler.RegisterUpdateUserHandler(mux, cfg)

	// chirp
	apihandler.RegisterCreateChirpHandler(mux, cfg)
	apihandler.RegisterGetAllChirpsHandler(mux, cfg)
	apihandler.RegisterGetChirpByIDHandler(mux, cfg)

	// auth
	apihandler.RegisterLoginHandler(mux, cfg)
	apihandler.RegisterRefreshHandler(mux, cfg)
	apihandler.RegisterRevokeHandler(mux, cfg)
}
