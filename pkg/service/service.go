package service

import (
	handlermodel "awesomeProject/pkg/handler/model"
	"awesomeProject/pkg/repository"
	servicemodel "awesomeProject/pkg/service/model"
)

type Transaction interface {
	AddTransaction(transaction handlermodel.Transaction) (*servicemodel.TransactionWithBalance, error)
	CancelLatestOddTransactions() error
}

type Service struct {
	Transaction
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Transaction: NewTransactionService(repos.Balance, repos.Transaction),
	}
}
