package main

import (
	idp "github.com/daschaa/open-idp/internal/idp"
	"github.com/daschaa/open-idp/internal/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	signingKey := []byte("your_secret_key")

	router := mux.NewRouter()
	server := idp.New(clientRepository, idp.WithSigningKey(&signingKey))
	router.HandleFunc("/token", server.TokenHandler)
	router.HandleFunc("/introspect", server.IntrospectHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
