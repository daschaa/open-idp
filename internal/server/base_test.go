package server_test

import (
	"github.com/stretchr/testify/suite"
	"open-idp/internal/repository"
	"testing"
	"time"
)

type testClock struct{}

func (testClock) Now() time.Time {
	return time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)
}

type serverSuite struct {
	suite.Suite
	clientRepository repository.ClientRepository
	clock            repository.Clock
	signingKey       []byte
}

func (suite *serverSuite) SetupTest() {
	suite.clock = testClock{}
	suite.clientRepository = repository.NewSimpleClientRepository()
	suite.signingKey = []byte("your_secret_key")
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(serverSuite))
}
