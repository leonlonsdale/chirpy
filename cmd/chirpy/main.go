package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/leonlonsdale/chirpy/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	_ = godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cfg := config{
		addr:           ":8080",
		FileserverHits: &atomic.Int32{},
		DBQueries:      *database.New(db),
		Platform:       os.Getenv("PLATFORM"),
		Secret:         os.Getenv("JWT_SECRET_KEY"),
		PolkaKey:       os.Getenv("POLKA_KEY"),
	}

	app := &application{
		config: cfg,
	}

	log.Fatal(app.run(app.mount()))

}
