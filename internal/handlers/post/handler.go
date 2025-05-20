package post_handler

import (
	post_client "pinstack-api-gateway/internal/clients/post"
	"pinstack-api-gateway/internal/logger"
)

type PostHandler struct {
	postClient post_client.PostClient
	log        *logger.Logger
}

func NewPostHandler(postClient post_client.PostClient, log *logger.Logger) *PostHandler {
	return &PostHandler{
		postClient: postClient,
		log:        log,
	}
}
