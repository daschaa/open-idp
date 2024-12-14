package main

import (
	"net/http"
	"open-idp/internal/repository"
	idp_server "open-idp/internal/server"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type SystemClock struct{}

func (c SystemClock) Now() time.Time {
	return time.Now()
}

func main() {
	clientRepository := repository.NewDynamoDbClientRepository(repository.NewDynamoDbClient())
	clock := SystemClock{}
	signingKey := []byte("your_secret_key")

	idp_server.InitIdpApi(clientRepository, clock, signingKey)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
