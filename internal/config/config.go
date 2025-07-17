package config

import (
	"net/http"
	"sync/atomic"

	"github.com/leonlonsdale/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBQueries      database.Queries
	Platform       string
	Secret         string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}
