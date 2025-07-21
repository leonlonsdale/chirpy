package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/util"
)

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
