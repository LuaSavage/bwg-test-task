package main

import (
	"context"
	"fmt"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/config"
	saccount "github.com/LuaSavage/bwg-test-task/service-b/internal/domain/adapter/storage/account"
	"github.com/LuaSavage/bwg-test-task/service-b/internal/domain/service/account"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/client/msgbroker/dto"
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

	ctx := context.TODO()
	// trying 5 times to establish connection
	pgxpool, err := postgresql.NewClient(ctx /*max retry attempts*/, 5, cfg.Storage)

	accountStorage := saccount.NewStorage(pgxpool)
	accountService := account.NewService(accountStorage, &logger)

	kafkaConsumer, err := msgbroker.NewKafkaConsumer(&dto.NewConsumerDTO{
		BrokerAdress:     fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port),
		GroupId:          "service_b_consumer_group",
		AutoOffsetReset:  "earliest",
		EnableAutoCommit: "false",
		Topic:            cfg.Kafka.Topic,
		Logger:           logger,
	})
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		err := kafkaConsumer.Close()
		logger.Error(err)
	}()

	consumerSerivce, err := consumer.NewConsumerService(kafkaConsumer, accountService, logger)
	consumerSerivce.Subscribe(ctx)
}
