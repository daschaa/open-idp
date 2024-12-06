package repository

import "time"

type Clock interface {
	Now() time.Time
}

type Client struct {
	ClientId     string
	ClientSecret string
}

type ClientRepository interface {
	SaveClient(clientId string, clientSecret string) (*Client, error)
	GetClient(clientId string) (*Client, error)
}
