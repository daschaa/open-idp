package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type client struct {
	ClientId     string `dynamodbav:"clientId"`
	ClientSecret string `dynamodbav:"clientSecret"`
}

type DynamoDbClientRepository struct {
	client *dynamodb.Client
}

func (r *DynamoDbClientRepository) SaveClient(clientId string, clientSecret string) (*Client, error) {
	client := client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	av, err := attributevalue.MarshalMap(client)
	if err != nil {
		return nil, err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("clients"),
		Item:      av,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
	}, nil
}

func (r *DynamoDbClientRepository) GetClient(clientId string) (*Client, error) {
	item, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("clients"),
		Key: map[string]types.AttributeValue{
			"clientId": &types.AttributeValueMemberS{Value: clientId},
		},
	})

	if err != nil {
		return nil, err
	}

	var client client

	err = attributevalue.UnmarshalMap(item.Item, &client)

	if err != nil {
		return nil, err
	}

	return &Client{
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
	}, nil
}

func NewDynamoDbClientRepository(client *dynamodb.Client) *DynamoDbClientRepository {
	return &DynamoDbClientRepository{
		client: client,
	}
}

func NewDynamoDbClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func NewLocalDynamoDbClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "DUMMYIDEXAMPLE", SecretAccessKey: "DUMMYEXAMPLEKEY", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}
