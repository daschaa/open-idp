package server

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"open-idp/internal/repository"
	"time"
)

type IntrospectRequest struct {
	Token *string `json:"token"`
}

type TokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

var signingKey = []byte("your_secret_key")

type IdpServer struct {
	clientRepository repository.ClientRepository
	clock            repository.Clock
}

func (s *IdpServer) ValidateClient(clientId string, clientSecret string) (bool, error) {
	client, err := s.clientRepository.GetClient(clientId)
	if err != nil {
		return false, err
	}

	if !(client.ClientSecret == clientSecret) {
		return false, errors.New("Secret does not match")
	}

	return true, nil
}

func (s *IdpServer) introspect(w http.ResponseWriter, r *http.Request) {
	request := IntrospectRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil || request.Token == nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(*request.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	active := claims.VerifyExpiresAt(s.clock.Now().Unix(), true)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"active": active})
}

func (s *IdpServer) token(w http.ResponseWriter, r *http.Request) {
	request := TokenRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if request.GrantType != "client_credentials" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	ok, err := s.ValidateClient(request.ClientId, request.ClientSecret)

	if !ok {
		http.Error(w, "Client is not authorized", http.StatusUnauthorized)
		return
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   request.ClientId,
		"exp":   s.clock.Now().Add(time.Hour).Unix(),
		"scope": "read:example",
	}).SignedString(signingKey)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"access_token": token, "token_type": "Bearer", "expires_in": "3600"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InitIdpApi(r repository.ClientRepository, clock repository.Clock) http.Handler {
	router := mux.NewRouter()
	server := IdpServer{
		clientRepository: r,
		clock:            clock,
	}
	router.HandleFunc("/token", server.token)
	router.HandleFunc("/introspect", server.introspect)
	return router
}
