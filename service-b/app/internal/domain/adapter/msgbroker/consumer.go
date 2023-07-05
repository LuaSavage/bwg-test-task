package msgbroker

import (
	"context"
	"encoding/json"

	tdto "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/dto"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/google/uuid"
)

type kafkaConsumer interface {
	ConsumeMessage() (response []byte)
}

type accountService interface {
	GetBalance(ctx context.Context, accountId uuid.UUID) (float64, error)
	Transfer(ctx context.Context, requestDTO *tdto.TransferRequestDTO) error
}

type ConsumerService struct {
	Consumer kafkaConsumer
	Service  accountService
	Logger   logging.Logger
}

func NewConsumerService(consumer kafkaConsumer, accountService accountService, logger logging.Logger) (*ConsumerService, error) {
	return &ConsumerService{
		Consumer: consumer,
		Service:  accountService,
		Logger:   logger,
	}, nil
}

func (c *ConsumerService) Subscribe(ctx context.Context) {
	for {
		msg := c.Consumer.ConsumeMessage()

		var reqDto tdto.TransferRequestDTO
		err := json.Unmarshal(msg, &reqDto)

		if err != nil {
			c.Logger.Error(err)
			continue
		}

		err = c.Service.Transfer(ctx, &reqDto)
		if err != nil {
			c.Logger.Error(err)
		}
	}
}
