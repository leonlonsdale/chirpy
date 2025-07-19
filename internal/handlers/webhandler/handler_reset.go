package webhandler

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func ResetHandler(db database.Queries, fs *atomic.Int32, platform string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if platform != "dev" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			log.Println("User accessed reset")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fs.Store(0)
		if err := db.ResetUsers(r.Context()); err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "there was a problem resetting users", err)
		}
		_, _ = fmt.Fprintf(w, "Hit counter reset!")
	}
}
