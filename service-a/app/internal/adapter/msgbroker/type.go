package msgbroker

import "github.com/google/uuid"

type KafkaTransferRequest struct {
	AccountID     uuid.UUID `json:"accountId"`
	TransactionId uuid.UUID `json:"transactionId"`
	Amount        float64   `json:"amount"`
}
