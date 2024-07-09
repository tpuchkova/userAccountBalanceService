package model

type Transaction struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transaction_id"`
}
