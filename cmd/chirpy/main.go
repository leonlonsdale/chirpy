package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/handlers"
	"github.com/leonlonsdale/chirpy/internal/storage"

	_ "github.com/lib/pq"
)

func main() {

	_ = godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	store := storage.NewStorage(db)

	cfg := config.Config{
		Addr:           os.Getenv("ADDR"),
		FileserverHits: &atomic.Int32{},
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("JWT_SECRET_KEY"),
		PolkaKey:       os.Getenv("POLKA_KEY"),
	}

	auth := auth.NewAuthService()
	handlers := handlers.NewHandlers(&store, &cfg, auth)

	app := &application{
		config:   cfg,
		store:    store,
		handlers: handlers,
		auth:     auth,
	}

	log.Fatal(app.run(app.mount()))

}
