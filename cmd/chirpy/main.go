package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
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
	cfg := &server.ApiConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      *database.New(db),
	}

	s := server.NewServer(port, cfg)
	fmt.Println("Starting server on 8080")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
