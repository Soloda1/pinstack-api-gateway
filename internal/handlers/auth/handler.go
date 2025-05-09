package auth_handler

import (
	auth_client "pinstack-api-gateway/internal/clients/auth"
	"pinstack-api-gateway/internal/logger"
)

type AuthHandler struct {
	authClient auth_client.AuthClient
	log        *logger.Logger
}

func NewAuthHandler(authClient auth_client.AuthClient, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		log:        log,
	}
}
