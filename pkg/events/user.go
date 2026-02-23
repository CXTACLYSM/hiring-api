package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	TopicUserCreated = "user.created"
)

type Event interface {
	Serialize() ([]byte, error)
	Topic() string
	Key() string
}

type UserCreated struct {
	UserId    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (uc *UserCreated) Serialize() ([]byte, error) {
	if uc == nil {
		return nil, errors.New("error serializing user created event: user created event is nil")
	}
	data, err := json.Marshal(uc)
	if err != nil {
		return nil, fmt.Errorf("error serializing user created event: %w", err)
	}

	return data, nil
}

func (uc *UserCreated) Topic() string {
	return TopicUserCreated
}

func (uc *UserCreated) Key() string {
	return uc.UserId
}
