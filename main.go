package main

import (
	"log"

	"github.com/justagabriel/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer("localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
