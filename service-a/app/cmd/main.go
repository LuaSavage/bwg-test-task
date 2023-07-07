package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	acchandler "github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/api/account"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/adapter/msgbroker"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/config"
	"github.com/LuaSavage/bwg-test-task/service-a/internal/domain/account"
	kafkaclient "github.com/LuaSavage/bwg-test-task/service-a/pkg/msgbroker"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg, err := config.GetConfig("", e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// initialising Kafka service
	kafkaProducer, err := kafkaclient.NewKafkaProducer(cfg.Kafka, e.Logger)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer kafkaProducer.Close()
	msgBroker := msgbroker.NewKafkaService(kafkaProducer, e.Logger)

	accountService := account.NewService(msgBroker)

	// initialising http handler's
	accountHandler := acchandler.NewHandler(accountService /*rate limit*/, cfg.Listen.RPC)
	accountHandler.Register(e)

	err = e.Start(fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receivedSignal := <-osSignal

	e.Logger.Infof("Application exit. Signal: %s (%d)", receivedSignal.String(), receivedSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
