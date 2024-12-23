package idp_test

import (
	"github.com/daschaa/open-idp/internal/idp"
	"github.com/daschaa/open-idp/internal/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type TestClock struct{}

func (TestClock) Now() time.Time {
	return time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)
}

type serverSuite struct {
	suite.Suite
	clientRepository repository.ClientRepository
	clock            repository.Clock
	signingKey       []byte
}

func (suite *serverSuite) SetupTest() {
	suite.clock = TestClock{}
	suite.clientRepository = repository.NewSimpleClientRepository()
	suite.signingKey = []byte("your_secret_key")
}

func (suite *serverSuite) InitIdpApi() http.Handler {
	router := mux.NewRouter()
	server := idp.New(suite.clientRepository, idp.WithSigningKey(&suite.signingKey), idp.WithClock(suite.clock))
	router.HandleFunc("/token", server.TokenHandler)
	router.HandleFunc("/introspect", server.IntrospectHandler)
	return router
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(serverSuite))
}
