package models

import (
	"time"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/relation/v1"
)

type Follower struct {
	ID         int64     `json:"id"`
	FollowerID int64     `json:"follower_id"`
	FolloweeID int64     `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type RelationUser struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func RelationUserFromProto(u *pb.User) *RelationUser {
	return &RelationUser{
		ID:        u.FollowerId,
		Username:  u.Username,
		AvatarURL: u.AvatarUrl,
	}
}
