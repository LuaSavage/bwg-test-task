package model

import "github.com/google/uuid"

type Account struct {
	Id      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}
