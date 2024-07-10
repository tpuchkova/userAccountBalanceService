package service

import (
	handlermodel "userAccountBalanceService/pkg/handler/model"
	"userAccountBalanceService/pkg/repository"
	servicemodel "userAccountBalanceService/pkg/service/model"
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
