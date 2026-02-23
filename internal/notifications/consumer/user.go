package consumer

import (
	"encoding/json"
	"log"

	"github.com/CXTACLYSM/hiring-api/pkg/events"
	"github.com/IBM/sarama"
)

type UserCreatedHandler struct{}

func NewUserCreatedHandler() *UserCreatedHandler {
	return &UserCreatedHandler{}
}

func (h *UserCreatedHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("consumer group session started, claims: %v", session.Claims())
	return nil
}

func (h *UserCreatedHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("consumer group session ended")
	return nil
}

func (h *UserCreatedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handleMessage(session, msg)
	}
	return nil
}

func (h *UserCreatedHandler) handleMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var event events.UserCreated
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("error unmarshaling message from topic=%s partition=%d offset=%d: %v",
			msg.Topic, msg.Partition, msg.Offset, err)
		session.MarkMessage(msg, "")
		return
	}

	log.Printf("received user.created: user_id=%s username=%s email=%s created_at=%s",
		event.UserId, event.Username, event.Email, event.CreatedAt.Format("2006-01-02 15:04:05"))

	session.MarkMessage(msg, "")
}
