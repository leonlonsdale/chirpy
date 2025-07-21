package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/storage"
	"github.com/leonlonsdale/chirpy/internal/util"
)

type application struct {
	config   config.Config
	store    storage.Storage
	handlers *handlers.Handlers
	auth     auth.Auth
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
	mux.Handle("POST /api/chirps", app.auth.JWTProtect(app.config.Secret, app.handlers.CreateChirp()))
	mux.Handle("GET /api/chirps", app.handlers.GetAllChirps())
	mux.Handle("GET /api/chirps/{chirpID}", app.handlers.GetChirpById())
	mux.Handle("DELETE /api/chirps/{chirpID}", app.auth.JWTProtect(app.config.Secret, app.handlers.DeleteChirpById()))

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

func (app *application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.config.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}

func (app *application) MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resString := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
			app.config.FileserverHits.Load())

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, resString)

	}
}

func (app *application) resetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.config.Platform != "dev" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Println("User accessed reset")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		app.config.FileserverHits.Store(0)
		if err := app.store.Users.Reset(r.Context()); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem resetting users", err)
		}
		_, _ = fmt.Fprintf(w, "Hit and Users counter reset!")
	}
}

func (app *application) createFileServer() http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))
}
