package dto

import "github.com/labstack/echo/v4"

type NewConsumerDTO struct {
	BrokerAdress     string
	GroupId          string
	AutoOffsetReset  string
	EnableAutoCommit string
	Topic            string
	Logger           echo.Logger
}
