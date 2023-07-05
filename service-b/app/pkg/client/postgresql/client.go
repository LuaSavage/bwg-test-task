package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func connectWithRetries(ctx context.Context, credentials string, maxRetries int, retryInterval time.Duration) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		pool, err = pgxpool.Connect(ctx, credentials)
		if err == nil {
			return pool, nil
		}
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("exceeded maximum postgresql connection retries")
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	credentials := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	pool, err = connectWithRetries(ctx, credentials, maxAttempts, time.Second*5)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
