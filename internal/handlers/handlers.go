package handlers

import (
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/storage"
)

type Handlers struct {
	// userHandlers  *UserHandlers
	ChirpHandlers *ChirpHandlers
}

func NewHandlers(store *storage.Storage, cfg *config.Config) *Handlers {
	return &Handlers{
		// userHandlers: &UserHandlers{
		// 	store: store,
		// 	cfg:   cfg,
		// },
		ChirpHandlers: &ChirpHandlers{
			Store: store,
			Cfg:   cfg,
		},
	}
}
