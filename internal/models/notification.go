package models

import (
	"encoding/json"
	"github.com/soloda1/pinstack-proto-definitions/events"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/notification/v1"
	"time"
)

// EventType used to swagger documentation and represents the type of event for a notification.
type EventType string

// NotificationSwagger represents the structure of a notification for Swagger documentation.
type NotificationSwagger struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Type      EventType       `json:"type"`
	IsRead    bool            `json:"is_read"`
	CreatedAt time.Time       `json:"created_at"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type Notification struct {
	ID        int64            `json:"id" db:"id"`
	UserID    int64            `json:"user_id" db:"user_id"`
	Type      events.EventType `json:"type" db:"type"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	Payload   json.RawMessage  `json:"payload,omitempty" db:"payload"`
}

func NotificationFromProto(notification *pb.NotificationResponse) *Notification {
	var createdAt time.Time
	if notification.CreatedAt != nil {
		createdAt = notification.CreatedAt.AsTime()
	}

	return &Notification{
		ID:        notification.Id,
		UserID:    notification.UserId,
		Type:      events.EventType(notification.Type),
		IsRead:    notification.IsRead,
		CreatedAt: createdAt,
		Payload:   notification.Payload,
	}
}
