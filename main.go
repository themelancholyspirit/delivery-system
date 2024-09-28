package main

import (
	"log"
)

func main() {

	storer, err := NewPostgreStorer()

	if err != nil {
		log.Fatalf("Failed to initialize storer: %v", err)
	}

	server := NewServer(":8080", storer)
	router := SetupRoutes(server)

	if err := router.Run(server.listenAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
