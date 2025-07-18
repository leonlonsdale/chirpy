package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/leonlonsdale/chirpy/internal/config"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/server"

	_ "github.com/lib/pq"
)

func main() {

	_ = godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	const port = "8080"
	cfg := &config.ApiConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      *database.New(db),
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("JWT_SECRET_KEY"),
		PolkaKey:       os.Getenv("POLKA_KEY"),
	}

	s := server.NewServer(port, cfg)
	server.RegisterHandlers(s.Mux, cfg)
	fmt.Println("Starting server on 8080")

	if err := s.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
