package main

import (
	"github.com/metamorfoso/recipe-book/server"
	"log"
)

func main() {
	host := "localhost"
	port := "8080" // TODO: pass in as env variable

	log.Println("Listening on port", port)
	err := server.RunServer(host, port)

	if err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
