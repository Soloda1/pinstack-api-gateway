package user_handler

import (
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/logger"
)

type UserHandler struct {
	userClient user_client.UserClient
	log        *logger.Logger
}

func NewUserHandler(userClient user_client.UserClient, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		log:        log,
	}
}
