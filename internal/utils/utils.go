package utils

import (
	"pinstack-api-gateway/internal/models"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func StringPtr(s string) *string {
	return &s
}

func GenerateUnknownAuthor() *models.User {
	return &models.User{
		ID:        0,
		Username:  "unknown",
		FullName:  StringPtr("Unknown Author"),
		AvatarURL: StringPtr("http://unknown.unknown"),
		Email:     "unknown@unknown.com",
		Bio:       StringPtr("Unknown Author BIO"),
	}
}
