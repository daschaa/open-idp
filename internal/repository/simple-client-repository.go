package repository

type SimpleClientRepository struct{}

type ClientNotFound struct {
	ClientId string
}

func (e ClientNotFound) Error() string {
	return "client not found"
}

func (r SimpleClientRepository) GetClient(clientId string) (*Client, error) {
	if clientId != "1234567890" {
		return nil, ClientNotFound{
			ClientId: clientId,
		}
	}
	return &Client{
		ClientId:     "1234567890",
		ClientSecret: "client_secret",
	}, nil
}

func (r SimpleClientRepository) SaveClient(clientId string, clientString string) (*Client, error) {
	return &Client{
		ClientId:     "1234567890",
		ClientSecret: "client_secret",
	}, nil
}

func NewSimpleClientRepository() SimpleClientRepository {
	return SimpleClientRepository{}
}
