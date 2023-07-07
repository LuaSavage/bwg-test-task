package msgbroker

import (
	"context"
	"encoding/json"

	accservice "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/service/account"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/google/uuid"
)

type kafkaConsumer interface {
	ConsumeMessage() (response []byte)
}

type accountService interface {
	GetBalance(ctx context.Context, accountId uuid.UUID) (float64, error)
	Transfer(ctx context.Context, requestDTO *accservice.TransferRequest) error
}

type ConsumerService struct {
	Consumer kafkaConsumer
	Service  accountService
	Logger   logging.Logger
	stopCh   chan struct{}
}

func NewConsumerService(consumer kafkaConsumer, accountService accountService, logger logging.Logger) (*ConsumerService, error) {
	return &ConsumerService{
		Consumer: consumer,
		Service:  accountService,
		Logger:   logger,
		stopCh:   make(chan struct{}),
	}, nil
}

func (c *ConsumerService) Subscribe(ctx context.Context) {
	for {
		select {
		case <-c.stopCh:
			return
		default:
			msg := c.Consumer.ConsumeMessage()

			var reqDto accservice.TransferRequest
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
}

func (c *ConsumerService) Close() {
	close(c.stopCh)
}
