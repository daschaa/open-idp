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
	clientRepository := repository.NewDynamoDbClientRepository(repository.NewDynamoDbClient())
	clock := SystemClock{}
	signingKey := []byte("your_secret_key")

	router := idp_server.InitIdpApi(clientRepository, clock, signingKey)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
