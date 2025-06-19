package relation_handler

import (
	relation_client "pinstack-api-gateway/internal/clients/relation"
	"pinstack-api-gateway/internal/logger"
)

type RelationHandler struct {
	relationClient relation_client.RelationClient
	log            *logger.Logger
}

func NewRelationHandler(relationClient relation_client.RelationClient, log *logger.Logger) *RelationHandler {
	return &RelationHandler{
		relationClient: relationClient,
		log:            log,
	}
}
