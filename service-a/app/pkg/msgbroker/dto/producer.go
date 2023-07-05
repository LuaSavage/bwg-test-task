package dto

import "github.com/labstack/echo/v4"

type NewProducerDTO struct {
	BrokerAdress string
	Topic        string
	Logger       echo.Logger
}
