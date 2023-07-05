package dto

import "github.com/google/uuid"

type KafkaTransferDTO struct {
	AccountID     uuid.UUID `json:"accountId"`
	TransactionId uuid.UUID `json:"transactionId"`
	Amount        float64   `json:"amount"`
}
