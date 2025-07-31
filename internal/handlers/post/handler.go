package post_handler

import (
	post_client "pinstack-api-gateway/internal/clients/post"
	user_client "pinstack-api-gateway/internal/clients/user"
	"pinstack-api-gateway/internal/logger"
)

type PostHandler struct {
	postClient post_client.PostClient
	userClient user_client.UserClient
	log        *logger.Logger
}

func NewPostHandler(postClient post_client.PostClient, userCLient user_client.UserClient, log *logger.Logger) *PostHandler {
	return &PostHandler{
		postClient: postClient,
		userClient: userCLient,
		log:        log,
	}
}
