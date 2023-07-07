package msgbroker

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type producer interface {
	ProduceMessage(messageKey []byte, messageValue []byte) error
	Close()
}

type KafkaService struct {
	Producer producer
	Logger   echo.Logger
}

func NewKafkaService(kafkaProducer producer, logger echo.Logger) *KafkaService {
	return &KafkaService{
		Producer: kafkaProducer,
		Logger:   logger,
	}
}

func (k *KafkaService) Transfer(ctx context.Context, dto *KafkaTransferRequest) error {
	payload, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	return k.Producer.ProduceMessage(nil, payload)
}
