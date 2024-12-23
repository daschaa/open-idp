package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	idp "github.com/daschaa/open-idp/internal/idp"
	repository "github.com/daschaa/open-idp/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting IDP idp")
	clientRepository := repository.NewDynamoDbClientRepository(repository.NewDynamoDbClient())

	router := mux.NewRouter()
	server := idp.New(clientRepository)
	router.HandleFunc("/token", server.TokenHandler)
	router.HandleFunc("/introspect", server.IntrospectHandler)

	lambda.Start(httpadapter.NewV2(router).ProxyWithContext)
}
