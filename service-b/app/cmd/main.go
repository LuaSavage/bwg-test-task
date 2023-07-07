package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/config"
	accstore "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/adapter/storage/account"
	"github.com/LuaSavage/bwg-test-task/service-b/internal/domain/service/account"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/client/postgresql"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"

	"github.com/LuaSavage/bwg-test-task/service-b/pkg/client/msgbroker"

	consumer "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/adapter/msgbroker"
)

func main() {
	logger := logging.GetLogger()
	cfg, err := config.GetConfig("")
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()
	// trying 5 times to establish connection
	pgxpool, err := postgresql.NewClient(ctx /*max retry attempts*/, 5, cfg.Storage)
	defer pgxpool.Close()
	if err != nil {
		logger.Fatal(err)
	}

	accountStorage := accstore.NewStorage(pgxpool)
	accountService := account.NewService(accountStorage, &logger)

	kafkaConsumer, err := msgbroker.NewKafkaConsumer(cfg.Kafka, logger)
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		err := kafkaConsumer.Close()
		logger.Error(err)
	}()

	consumerSerivce, err := consumer.NewConsumerService(kafkaConsumer, accountService, logger)
	if err != nil {
		logger.Fatal(err)
	}
	go consumerSerivce.Subscribe(ctx)

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receivedSignal := <-osSignal
	logger.Infof("Application exit. Signal: %s (%d)", receivedSignal.String(), receivedSignal)
	consumerSerivce.Close()
}
