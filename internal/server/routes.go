package server

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers/apihandler"
	"github.com/leonlonsdale/chirpy/internal/handlers/webhandler"
)

func RegisterHandlers(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /admin/metrics", webhandler.MetricsHandler(cfg))
	mux.HandleFunc("POST /admin/reset", webhandler.ResetHandler(cfg))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))))

	// ==========[ API ]
	mux.HandleFunc("GET /api/healthz", apihandler.HealthzHandler)

	// user
	mux.HandleFunc("POST /api/users", apihandler.CreateUserHandler(cfg))
	mux.HandleFunc("PUT /api/users", apihandler.UpdateUserHandler(cfg))

	// chirp
	mux.Handle("POST /api/chirps", auth.MiddlewareCheckJWT(cfg.Secret, apihandler.CreateChirpHandler(cfg)))
	mux.HandleFunc("GET /api/chirps", apihandler.GetAllChirpsHandler(cfg))
	mux.HandleFunc("GET /api/chirps/{chirpID}", apihandler.GetChirpByIDHandler(cfg))
	mux.Handle("DELETE /api/chirps/{chirpID}", auth.MiddlewareCheckJWT(cfg.Secret, apihandler.DeleteChirpByID(cfg)))

	// auth
	mux.HandleFunc("POST /api/login", apihandler.LoginHandler(cfg))
	mux.HandleFunc("POST /api/refresh", apihandler.RefreshHandler(cfg))
	mux.HandleFunc("POST /api/revoke", apihandler.RevokeHandler(cfg))
}
