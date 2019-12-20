package testtask1

import "github.com/shopspring/decimal"

type TransactionRequest struct {
	State         StateType       `json:"state"`
	Amount        decimal.Decimal `json:"amount,string"`
	TransactionId string          `json:"transactionId"`
}
type TransactionResponse struct {
	Error   string          `json:"error"`
	Status  StatusType      `json:"status"`
	Balance decimal.Decimal `json:"balance"`
}
