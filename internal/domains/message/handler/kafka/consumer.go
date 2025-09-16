package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"

	"MessageService/internal/domains/message/model"
	"MessageService/internal/domains/message/usecase"
)

type MessageConsumer struct {
	reader  *kafka.Reader
	useCase usecase.MessageUseCase
	ctx     context.Context
}

func NewMessageConsumer(brokers []string, topic, groupID string, uc usecase.MessageUseCase) *MessageConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &MessageConsumer{
		reader:  r,
		useCase: uc,
		ctx:     context.Background(),
	}
}

func (c *MessageConsumer) Start() error {
	for {
		m, err := c.reader.ReadMessage(c.ctx)
		if err != nil {
			log.Printf("[Kafka] Error reading message: %v\n", err)
			continue
		}

		var dto MessageDTO
		if err := json.Unmarshal(m.Value, &dto); err != nil {
			log.Printf("[Kafka] JSON unmarshal error: %v\n", err)
			continue
		}

		msg := dtoToModel(dto)

		if err := c.useCase.Save(c.ctx, msg); err != nil {
			log.Printf("[Kafka] Save error: %v\n", err)
		} else {
			log.Printf("[Kafka] Message saved: %+v\n", msg)
		}
	}
}

func (c *MessageConsumer) Close() error {
	return c.reader.Close()
}

func dtoToModel(dto MessageDTO) *model.Message {
	return &model.Message{
		T:  dto.T,
		ST: dto.ST,
		Pos: model.Pos{
			X:  dto.Pos.X,
			Y:  dto.Pos.Y,
			Z:  dto.Pos.Z,
			A:  dto.Pos.A,
			S:  dto.Pos.S,
			St: dto.Pos.St,
		},
		Params: dto.Params,
	}
}
