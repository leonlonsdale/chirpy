package config

import (
	"sync/atomic"
)

type Config struct {
	Addr           string
	FileserverHits *atomic.Int32
	Platform       string
	Secret         string
	PolkaKey       string
}
