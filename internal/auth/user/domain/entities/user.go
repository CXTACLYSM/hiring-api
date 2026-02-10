package entities

import "time"

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
