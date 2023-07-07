package account

import "github.com/google/uuid"

type TransferRequest struct {
	AccountID     uuid.UUID `json:"accountId"`
	TransactionId uuid.UUID `json:"transactionId"`
	Amount        float64   `json:"amount"`
}
