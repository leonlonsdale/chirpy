package main

import (
	"fmt"
	"log"

	"github.com/ionztorm/chirpy/internal/server"
)

func main() {
	const port = "8080"

	s := server.NewServer(port)
	fmt.Println("Starting server on 8080")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
