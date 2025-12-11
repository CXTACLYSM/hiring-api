package entity

import "time"

type User struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Age          uint8  `json:"age"`
	Gender       string `json:"gender"`
	Language     string `json:"language"`
	Type         string `json:"type"`
	Status       string `json:"status"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
