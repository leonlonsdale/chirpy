package main

import (
	"log"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers/apihandler"
	"github.com/leonlonsdale/chirpy/internal/handlers/webhandler"
	"github.com/leonlonsdale/chirpy/internal/storage"
)

type application struct {
	config config.Config
	store  storage.Storage
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /admin/metrics", webhandler.MetricsHandler(app.config.DBQueries, app.config.FileserverHits))
	mux.Handle("POST /admin/reset", webhandler.ResetHandler(app.config.DBQueries, app.config.FileserverHits, app.config.Platform))
	mux.Handle("/app/", app.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))))

	// ==========[ API ]
	mux.HandleFunc("GET /api/healthz", app.healthHandler)

	// user
	mux.Handle("POST /api/users", apihandler.CreateUserHandler(app.config.DBQueries))
	mux.Handle("PUT /api/users", apihandler.UpdateUserHandler(app.config.DBQueries, app.config.Secret))

	// chirp
	mux.Handle("POST /api/chirps", auth.MiddlewareCheckJWT(app.config.Secret, apihandler.CreateChirpHandler(app.config.DBQueries)))
	mux.Handle("GET /api/chirps", apihandler.GetAllChirpsHandler(app.config.DBQueries))
	mux.Handle("GET /api/chirps/{chirpID}", apihandler.GetChirpByIDHandler(app.config.DBQueries))
	mux.Handle("DELETE /api/chirps/{chirpID}", auth.MiddlewareCheckJWT(app.config.Secret, apihandler.DeleteChirpByID(app.config.DBQueries)))

	// auth
	mux.Handle("POST /api/login", apihandler.LoginHandler(app.config.DBQueries, app.config.Secret))
	mux.Handle("POST /api/refresh", apihandler.RefreshHandler(app.config.DBQueries, app.config.Secret))
	mux.Handle("POST /api/revoke", apihandler.RevokeHandler(app.config.DBQueries))

	// webhooks
	mux.Handle("POST /api/polka/webhooks", apihandler.UpgradeToChirpyRedHandler(app.config.DBQueries, app.config.PolkaKey))

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

func (app *application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.config.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}
