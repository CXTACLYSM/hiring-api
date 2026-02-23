package kafka

import (
	"fmt"
	"time"

	"github.com/CXTACLYSM/hiring-api/pkg/events"
	"github.com/IBM/sarama"
)

type Publisher struct {
	producer sarama.SyncProducer
}

func NewPublisher(producer sarama.SyncProducer) *Publisher {
	return &Publisher{
		producer: producer,
	}
}

func (p *Publisher) Push(event events.Event) error {
	serialized, err := event.Serialize()
	if err != nil {
		return fmt.Errorf("error pushing event to kafka: %w", err)
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: event.Topic(),
		Key:   sarama.StringEncoder(event.Key()),
		Value: sarama.ByteEncoder(serialized),
		//Headers:   nil,
		//Metadata:  nil,
		//Offset:    0,
		//Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error pushing event to kafka: %w", err)
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.producer.Close()
}
