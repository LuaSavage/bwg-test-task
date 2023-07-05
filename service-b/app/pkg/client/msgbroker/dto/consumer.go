package dto

import (
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
)

type NewConsumerDTO struct {
	BrokerAdress     string
	GroupId          string
	AutoOffsetReset  string
	EnableAutoCommit string
	Topic            string
	Logger           logging.Logger
}
