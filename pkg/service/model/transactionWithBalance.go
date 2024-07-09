package model

type TransactionWithBalance struct {
	ID            int
	State         string
	Amount        float64
	TransactionID string
	UserBalance   float64
}
