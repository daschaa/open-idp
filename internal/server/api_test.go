package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	idp_server "open-idp/internal/server"
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
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then the response should be 200 OK with an access token
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "application/json", response.Header().Get("Content-Type"))
	var jsonTokenResponse tokenResponse
	err := json.NewDecoder(response.Body).Decode(&jsonTokenResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tokenResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.1bY6cMhkY-h6dATMYi6-KmTlPiWO-DlIbHrHONXQbQs",
		ExpiresIn:   "3600",
		TokenType:   "Bearer",
	}, jsonTokenResponse)
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsBadRequestIfBodyCanNotBeParsed() {
	// when sending a POST request to /token with an invalid body
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(`clientId: 1234567890`))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then the response should be 400 Bad TokenRequest
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Invalid body\n", response.Body.String())
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsBadRequestIfGrantTypeWrong() {
	// when sending a POST request to /token without the correct grant type
	requestBody := `{"client_id":"1234567890","client_secret":"client_secret", "grant_type":"password"}`
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then the response should be 400 Bad TokenRequest with an error message
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Unsupported grant type\n", response.Body.String())
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsUnauthorizedIfClientSecretIsWrong() {
	// when sending a POST request to /token with the wrong client secret
	requestBody := `{"client_id":"1234567890","client_secret":"wrong_secret","grant_type":"client_credentials"}`

	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then the response should be 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, response.Result().StatusCode)
}

func (suite *serverSuite) Test_TokenEndpoint_ReturnsUnauthorizedIfClientNotInDatabase() {
	// when sending a POST request to /token
	requestBody := `{"client_id":"123","client_secret":"client_secret","grant_type":"client_credentials"}`
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then the response should be 200 OK with an access token
	assert.Equal(suite.T(), http.StatusUnauthorized, response.Result().StatusCode)
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsActive() {
	// when sending a POST request with valid token to /introspect
	requestBody := `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.1bY6cMhkY-h6dATMYi6-KmTlPiWO-DlIbHrHONXQbQs"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":true}\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsBadRequestIfInvalidBody() {
	// when sending a POST request with valid token to /introspect
	requestBody := `{"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTE1MjQwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.1bY6cMhkY-h6dATMYi6-KmTlPiWO-DlIbHrHONXQbQs"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusBadRequest, response.Result().StatusCode)
	assert.Equal(suite.T(), "Invalid body\n", response.Body.String())
}

func (suite *serverSuite) Test_IntrospectEndpoint_ReturnsIsNotActiveIfExpired() {
	// when sending a POST request with valid token to /introspect
	requestBody := `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkxNTEwMDAwMCwic2NvcGUiOiJyZWFkOmV4YW1wbGUiLCJzdWIiOiIxMjM0NTY3ODkwIn0.G0mBIfxTwyQ50CuKBn3Ti5Dj_w8PsktulHq9R-qpLzQ"}`
	request := httptest.NewRequest(http.MethodPost, "/introspect", bytes.NewBufferString(requestBody))
	response := httptest.NewRecorder()
	idp_server.InitIdpApi(suite.clientRepository, suite.clock).ServeHTTP(response, request)

	// then its returned as active
	assert.Equal(suite.T(), http.StatusOK, response.Result().StatusCode)
	assert.Equal(suite.T(), "{\"active\":false}\n", response.Body.String())
}
