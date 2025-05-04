package models

import (
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserFromProto(u *pb.User) *User {
	return &User{
		ID:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Bio:       u.Bio,
		AvatarURL: u.AvatarUrl,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Bio:       u.Bio,
		AvatarUrl: u.AvatarURL,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
