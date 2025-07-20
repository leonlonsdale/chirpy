package config

import (
	"sync/atomic"

	"github.com/leonlonsdale/chirpy/internal/database"
)

type Config struct {
	Addr           string
	FileserverHits *atomic.Int32
	DBQueries      database.Queries
	Platform       string
	Secret         string
	PolkaKey       string
}
