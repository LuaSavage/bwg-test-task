package main

import (
	"fmt"

	haccount "github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/api/account"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/msgbroker"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/config"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/domain/account"
	"github.com/LuaSavage/bwg-test-task/service-a/pkg/msgbroker/dto"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg, err := config.GetConfig("", e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// initialising Kafka service
	producerCfg := &dto.NewProducerDTO{
		BrokerAdress: fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port),
		Topic:        cfg.Kafka.Topic,
		Logger:       e.Logger,
	}
	msgBroker, err := msgbroker.NewKafkaService(producerCfg)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// initialising http handler's
	accountService := account.NewService(msgBroker)
	accountHandler := haccount.NewHandler(accountService /*rate limit*/, cfg.Listen.RPC)
	accountHandler.Register(e)
}
