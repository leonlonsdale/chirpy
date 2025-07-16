package server

import (
	"net/http"
	"sync/atomic"

	"github.com/ionztorm/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBQueries      database.Queries
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}
