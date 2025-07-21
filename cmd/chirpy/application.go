package main

import (
	"log"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/storage"
)

type application struct {
	config   *config.Config
	store    *storage.Storage
	handlers *handlers.Handlers
	auth     *auth.Auth
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/app/", app.MiddlewareMetricsInc(app.createFileServer()))
	mux.Handle("GET /admin/metrics", app.MetricsHandler())
	mux.Handle("POST /admin/reset", app.resetHandler())
	mux.Handle("GET /api/healthz", app.healthHandler())

	// user
	mux.Handle("POST /api/users", app.handlers.CreateUser())
	mux.Handle("PUT /api/users", app.handlers.UpdateUser())

	// chirp
	mux.Handle("POST /api/chirps", app.auth.JWTProtect(app.handlers.CreateChirp()))
	mux.Handle("GET /api/chirps", app.handlers.GetAllChirps())
	mux.Handle("GET /api/chirps/{chirpID}", app.handlers.GetChirpById())
	mux.Handle("DELETE /api/chirps/{chirpID}", app.auth.JWTProtect(app.handlers.DeleteChirpById()))

	// auth
	mux.Handle("POST /api/login", app.handlers.Login())
	mux.Handle("POST /api/refresh", app.handlers.Refresh())
	mux.Handle("POST /api/revoke", app.handlers.Revoke())

	// webhooks
	mux.Handle("POST /api/pokja/webhooks", app.handlers.UpdateUser())

	return mux
}

func (app *application) run(mux *http.ServeMux) error {

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server listening on %s", app.config.Addr)

	return srv.ListenAndServe()
}
