package repository

import (
	"userAccountBalanceService/pkg/repository/model"

	"github.com/jmoiron/sqlx"
)

type Balance interface {
	SaveBalance(balance float64) error
	GetBalance() (float64, error)
}

type Transaction interface {
	AddTransaction(transaction model.Transaction) (model.Transaction, error)
	GetLatestTransactions(count int) ([]model.Transaction, error)
	CancelByIDs(ids []int) error
}

type Repository struct {
	Balance
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Balance:     NewBalancePostgres(db),
		Transaction: NewTransactionPostgres(db),
	}
}
