package idp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/daschaa/open-idp/internal/idp"
	"github.com/daschaa/open-idp/internal/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsAccessToken() {
	// when sending a POST request to /token
	requestBody := `{"client_id":"1234567890","client_secret":"client_secret","grant_type":"client_credentials"}`
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then the response should be 200 OK with an access TokenHandler
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "application/json", response.Header().Get("Content-Type"))
	var jsonTokenResponse tokenResponse
	err := json.NewDecoder(response.Body).Decode(&jsonTokenResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tokenResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEwNDEzNzk1NjAwLCJzY29wZSI6InJlYWQ6ZXhhbXBsZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.SEe81GxuVobvhb4dWVb7L7gvlZRXrDzs99riasc-OmA",
		ExpiresIn:   "3600",
		TokenType:   "Bearer",
	}, jsonTokenResponse)
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsBadRequestIfBodyCanNotBeParsed() {
	// when sending a POST request to /token with an invalid body
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(`clientId: 1234567890`))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then the response should be 400 Bad tokenRequest
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Invalid body\n", response.Body.String())
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsBadRequestIfGrantTypeWrong() {
	// when sending a POST request to /token without the correct grant type
	requestBody := `{"client_id":"1234567890","client_secret":"client_secret", "grant_type":"password"}`
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then the response should be 400 Bad tokenRequest with an error message
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Unsupported grant type\n", response.Body.String())
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsUnauthorizedIfClientSecretIsWrong() {
	// when sending a POST request to /token with the wrong client secret
	requestBody := `{"client_id":"1234567890","client_secret":"wrong_secret","grant_type":"client_credentials"}`

	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then the response should be 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, response.Result().StatusCode)
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsUnauthorizedIfClientNotInDatabase() {
	// when sending a POST request to /token
	requestBody := `{"client_id":"123","client_secret":"client_secret","grant_type":"client_credentials"}`
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then the response should be 200 OK with an access TokenHandler
	assert.Equal(suite.T(), http.StatusUnauthorized, response.Result().StatusCode)
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsActive() {
	// when sending a POST request with valid TokenHandler to /introspect
	requestBody := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEwNDEzNzk1NjAwLCJzY29wZSI6InJlYWQ6ZXhhbXBsZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.SEe81GxuVobvhb4dWVb7L7gvlZRXrDzs99riasc-OmA"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":true,\"sub\":\"1234567890\"}\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsBadRequestIfInvalidBody() {
	// when sending a POST request with valid TokenHandler to /introspect
	requestBody := `{"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.1bY6cMhkY-h6dATMYi6-KmTlPiWO-DlIbHrHONXQbQs"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Invalid body\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsNotActiveIfExpired() {
	// when sending a POST request with valid TokenHandler to /introspect
	requestBody := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTEwMDAwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.G0mBIfxTwyQ50CuKBn3Ti5Dj_w8PsktulHq9R-qpLzQ"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":false}\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsNotActiveIfClientIsUnknown() {
	// when sending a POST request with valid TokenHandler to /introspect
	requestBody := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiJzb21lLWNsaWVudCJ9.lNBsxlRqA-NpYw_TY_tO_q5cVAb3zWDXMCerh3suxP0"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":false}\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsNotActiveIfSignatureInvalid() {
	// when sending a POST request with valid TokenHandler to /introspect
	requestBody := `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.Xd3ZGLgnkvhJiIJqBR937N0v-_PWPEotAWZPo2qWabE"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	suite.InitIdpApi().ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":false}\n", response.Body.String())
}

// ExampleTokenHandler demonstrates how to use the TokenHandler to generate a new token.
func ExampleTokenHandler() {
	signingKey := []byte("test")
	server := idp.New(repository.NewSimpleClientRepository(), idp.WithClock(TestClock{}), idp.WithSigningKey(&signingKey))
	requestBody := `{"client_id":"1234567890","client_secret":"client_secret","grant_type":"client_credentials"}`

	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	server.TokenHandler(response, request)

	fmt.Println(response.Result().StatusCode)
	fmt.Println(response.Header().Get("Content-Type"))
	fmt.Println(response.Body.String())
	// Output:
	// 200
	// application/json
	// {"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEwNDEzNzk1NjAwLCJzY29wZSI6InJlYWQ6ZXhhbXBsZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.ujFvivtthZOmmKe_BlxhdMNVh6UHAU5bYww8y62OmTI","expires_in":"3600","token_type":"Bearer"}
}

// ExampleIntrospectHandler demonstrates how to use the IntrospectHandler to introspect a token.
func ExampleIntrospectHandler() {
	signingKey := []byte("test")
	server := idp.New(repository.NewSimpleClientRepository(), idp.WithClock(TestClock{}), idp.WithSigningKey(&signingKey))
	requestBody := `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEwNDEzNzk1NjAwLCJzY29wZSI6InJlYWQ6ZXhhbXBsZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.ujFvivtthZOmmKe_BlxhdMNVh6UHAU5bYww8y62OmTI"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	server.IntrospectHandler(response, request)

	fmt.Println(response.Result().StatusCode)
	fmt.Println(response.Header().Get("Content-Type"))
	fmt.Println(response.Body.String())
	// Output:
	// 200
	// application/json
	// {"active":true,"sub":"1234567890"}
}
