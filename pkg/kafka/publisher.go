package kafka

import (
	"fmt"
	"io"

	"github.com/CXTACLYSM/hiring-api/pkg/events"
	"github.com/IBM/sarama"
)

type EventPublisher interface {
	io.Closer
	Push(event events.Event) error
}

type DefaultEventPublisher struct {
	producer sarama.SyncProducer
}

func NewDefaultEventPublisher(producer sarama.SyncProducer) *DefaultEventPublisher {
	return &DefaultEventPublisher{
		producer: producer,
	}
}

func (p *DefaultEventPublisher) Push(event events.Event) error {
	serialized, err := event.Serialize()
	if err != nil {
		return fmt.Errorf("error pushing event to kafka: %w", err)
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: event.Topic(),
		Key:   sarama.StringEncoder(event.Key()),
		Value: sarama.ByteEncoder(serialized),
	})
	if err != nil {
		return fmt.Errorf("error pushing event to kafka: %w", err)
	}

	return nil
}

func (p *DefaultEventPublisher) Close() error {
	return p.producer.Close()
}
