package handlers

import (
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/storage"
)

type Users interface {
	CreateUser() http.HandlerFunc
	UpdateUser() http.HandlerFunc
}

type Chirps interface {
	CreateChirp() http.HandlerFunc
	GetAllChirps() http.HandlerFunc
	GetChirpById() http.HandlerFunc
	DeleteChirpById() http.HandlerFunc
}

type Auth interface {
	Login() http.HandlerFunc
	Refresh() http.HandlerFunc
	Revoke() http.HandlerFunc
}

type Webhooks interface {
	UpgradeUser() http.HandlerFunc
}

type Handlers struct {
	Users
	Chirps
	Auth
	Webhooks
}

func NewHandlers(store *storage.Storage, cfg *config.Config, auth *auth.Auth) *Handlers {
	return &Handlers{
		Users: &UserHandlers{
			store: store,
			cfg:   cfg,
			auth:  auth,
		},
		Chirps: &ChirpHandlers{
			store: store,
			cfg:   cfg,
		},
		Auth: &AuthHandlers{
			store: store,
			cfg:   cfg,
			auth:  auth,
		},
		Webhooks: &WebhookHandlers{
			store: store,
			cfg:   cfg,
			auth:  auth,
		},
	}
}
