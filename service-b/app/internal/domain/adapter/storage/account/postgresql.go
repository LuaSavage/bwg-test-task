package account

import (
	"context"
	"fmt"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/domain/model"
	saccount "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/service/account"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/client/postgresql"
	"github.com/jackc/pgx"

	"github.com/google/uuid"
)

type Storage struct {
	client postgresql.Client
	//logger *logging.Logger
}

func NewStorage(client postgresql.Client /*logger *logging.Logger*/) *Storage {
	return &Storage{
		client: client,
		//logger: logger,
	}
}

func (s *Storage) GetBalance(ctx context.Context, accountId uuid.UUID) (*model.Account, error) {
	query := `SELECT balance from account WHERE id = $1;`
	row := s.client.QueryRow(ctx, query, accountId)
	var account model.Account
	err := row.Scan(&account.Balance)
	account.Id = accountId
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account by id = %s was not found", accountId.String())
		}
		return nil, err
	}
	return &account, nil
}

func (s *Storage) Withdraw(ctx context.Context, accountId uuid.UUID, amount float64) (saccount.Transaction, error) {
	tx, err := s.client.Begin(ctx)
	if err != nil {
		return nil, err
	}

	query := `UPDATE account SET balance = balance - $1 WHERE id = $2 AND balance >= $1 RETURNING * FOR UPDATE;`
	row := tx.QueryRow(ctx, query, amount, accountId)

	var account model.Account
	err = row.Scan(&account.Id, &account.Balance)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, err
	}

	return tx, nil
}
