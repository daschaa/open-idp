package integrationtest

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/daschaa/open-idp/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type dynamoDbSuite struct {
	suite.Suite
	repository *repository.DynamoDbClientRepository
	client     *dynamodb.Client
}

func (suite *dynamoDbSuite) SetupTest() {
	suite.repository = repository.NewDynamoDbClientRepository(repository.NewLocalDynamoDbClient())
	suite.client = repository.NewLocalDynamoDbClient()
}

func TestDynamoDbSuite(t *testing.T) {
	suite.Run(t, new(dynamoDbSuite))
}

func (s *dynamoDbSuite) Test_DynamoDbClientRepository_SaveClient() {
	// when saving a client
	savedClient, err := s.repository.SaveClient("123456789", "client_secret")

	// then the client should be saved
	assert.NoError(s.T(), err)
	item, err := s.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("clients"),
		Key: map[string]types.AttributeValue{
			"clientId": &types.AttributeValueMemberS{Value: "123456789"},
		},
	})
	assert.NoError(s.T(), err)
	client := &repository.Client{}
	err = attributevalue.UnmarshalMap(item.Item, client)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "123456789", client.ClientId)
	assert.Equal(s.T(), "client_secret", client.ClientSecret)
	assert.Equal(s.T(), savedClient, client)
}

func (s *dynamoDbSuite) Test_DynamoDbClientRepository_GetClient() {
	// when getting a saved client
	savedClient, err := s.repository.SaveClient("123456789", "client_secret")
	client, err := s.repository.GetClient("123456789")

	// then the client should be valid
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "123456789", client.ClientId)
	assert.Equal(s.T(), "client_secret", client.ClientSecret)
	assert.Equal(s.T(), client, savedClient)
}
