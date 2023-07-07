package account

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/domain/model"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/google/uuid"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type storage interface {
	GetBalance(ctx context.Context, accountId uuid.UUID) (*model.Account, error)
	Withdraw(ctx context.Context, accountId uuid.UUID, amount float64) (Transaction, error)
}

type Service struct {
	storage storage
	logger  *logging.Logger
}

func NewService(storage storage, logger *logging.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) GetBalance(ctx context.Context, accountId uuid.UUID) (float64, error) {
	s.logger.Infof("Getting balance by account id = %s", accountId.String())
	acc, err := s.storage.GetBalance(ctx, accountId)
	return acc.Balance, err
}

func (s *Service) Transfer(ctx context.Context, request *TransferRequest) error {
	if request.Amount <= 0 {
		return fmt.Errorf("funds amount must be positive")
	}

	balance, err := s.GetBalance(ctx, request.AccountID)
	s.logger.Infof("Starting funds transer from accId = %s, amount = %f, transactionId = %s",
		request.AccountID, request.Amount, request.TransactionId.String())

	if err != nil {
		return err
	}

	if balance < request.Amount {
		return fmt.Errorf("not enough money")
	}

	tx, err := s.storage.Withdraw(ctx, request.AccountID, request.Amount)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	if rand.Intn(2) == 0 {
		s.logger.Infof("Transfer transaction has failed accId = %s, amount = %f, transactionId = %s",
			request.AccountID, request.Amount, request.TransactionId.String())
		err = tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("transaction error: %s", err.Error())
		}
		return nil
	}

	s.logger.Infof("Transfer transaction has succeeded accId = %s, amount = %f, transactionId = %s",
		request.AccountID, request.Amount, request.TransactionId.String())
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction error: %s", err.Error())
	}

	return nil
}
