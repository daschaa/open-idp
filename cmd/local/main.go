package main

import (
	"log"
	"net/http"
	"open-idp/internal/repository"
	idp_server "open-idp/internal/server"
	"time"
)

type SystemClock struct{}

func (c SystemClock) Now() time.Time {
	return time.Now()
}

func main() {
	clientRepository := repository.NewDynamoDbClientRepository(repository.NewLocalDynamoDbClient())
	_, err := clientRepository.SaveClient("1234567890", "client_secret")
	if err != nil {
		log.Fatalf("Failed to save client: %v", err)
	}

	clock := SystemClock{}

	router := idp_server.InitIdpApi(clientRepository, clock)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
