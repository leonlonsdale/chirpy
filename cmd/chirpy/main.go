package main

import (
	"fmt"
	"log"

	"github.com/ionztorm/chirpy/internal/server"
)

func main() {
	s := server.NewServer()
	fmt.Println("Starting server on 8080")

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
