package account

import (
	"context"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/api/dto"
	kdto "github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/msgbroker/dto"
)

type msgBroker interface {
	Transfer(ctx context.Context, requestDTO *kdto.KafkaTransferDTO) error
}

type Service struct {
	msgbroker msgBroker
}

func NewService(m msgBroker) *Service {
	return &Service{
		msgbroker: m,
	}
}

func (s *Service) Transfer(ctx context.Context, requestDTO *dto.TransferRequestDTO) error {

	/*
		some business logic here
	*/

	return s.msgbroker.Transfer(ctx, &kdto.KafkaTransferDTO{
		AccountID:     requestDTO.AccountID,
		TransactionId: requestDTO.TransactionId,
		Amount:        requestDTO.Amount,
	})
}
