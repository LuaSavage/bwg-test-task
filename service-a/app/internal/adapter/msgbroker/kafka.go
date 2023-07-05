package msgbroker

import (
	"context"
	"encoding/json"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/msgbroker/dto"
	kafkaclient "github.com/LuaSavage/bwg-test-task/service-a/pkg/msgbroker"
	kdto "github.com/LuaSavage/bwg-test-task/service-a/pkg/msgbroker/dto"
	"github.com/labstack/echo/v4"
)

type KafkaService struct {
	Producer *kafkaclient.KafkaProducer
	Logger   echo.Logger
}

func NewKafkaService(dto *kdto.NewProducerDTO) (*KafkaService, error) {
	kafkaProducer, err := kafkaclient.NewKafkaProducer(dto)
	if err != nil {
		return nil, err
	}

	return &KafkaService{
		Producer: kafkaProducer,
		Logger:   dto.Logger,
	}, nil
}

func (k *KafkaService) Transfer(ctx context.Context, dto *dto.KafkaTransferDTO) error {
	payload, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	return k.Producer.ProduceMessage(nil, payload)
}
