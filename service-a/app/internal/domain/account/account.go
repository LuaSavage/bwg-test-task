package account

import (
	"context"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/msgbroker"
)

type msgBroker interface {
	Transfer(ctx context.Context, request *msgbroker.KafkaTransferRequest) error
}

type Service struct {
	msgbroker msgBroker
}

func NewService(m msgBroker) *Service {
	return &Service{
		msgbroker: m,
	}
}

func (s *Service) Transfer(ctx context.Context, request *TransferRequest) error {
	return s.msgbroker.Transfer(ctx, &msgbroker.KafkaTransferRequest{
		AccountID:     request.AccountID,
		TransactionId: request.TransactionId,
		Amount:        request.Amount,
	})
}
