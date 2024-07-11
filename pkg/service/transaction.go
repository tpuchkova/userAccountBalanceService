package service

import (
	"errors"
	"log"
	"strconv"

	handlermodel "userAccountBalanceService/pkg/handler/model"
	"userAccountBalanceService/pkg/repository"
	repositorymodel "userAccountBalanceService/pkg/repository/model"
	servicemodel "userAccountBalanceService/pkg/service/model"
)

const (
	stateWin = "win"
)

type TransactionService struct {
	balanceRepo     repository.Balance
	transactionRepo repository.Transaction
}

func NewTransactionService(balanceRepo repository.Balance, transactionRepo repository.Transaction) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
	}
}

func (s *TransactionService) AddTransaction(transaction handlermodel.Transaction) (*servicemodel.TransactionWithBalance, error) {
	amount, err := strconv.ParseFloat(transaction.Amount, 2)
	if err != nil {
		return nil, errors.New("invalid amount")
	}

	balance, err := s.balanceRepo.GetBalance()
	if err != nil {
		return nil, err
	}

	newBalance := calculateBalance(transaction.State, balance, amount)

	if newBalance < 0 {
		return nil, errors.New("insufficient balance")
	}

	t, err := s.transactionRepo.AddTransaction(repositorymodel.Transaction{
		State:         transaction.State,
		Amount:        amount,
		TransactionID: transaction.TransactionID,
	})
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"transactions_transaction_id_key\"" {
			return nil, errors.New("transaction already exists")
		}
		return nil, err
	}

	if err := s.balanceRepo.SaveBalance(newBalance); err != nil {
		return nil, err
	}

	return &servicemodel.TransactionWithBalance{
		ID:            t.ID,
		State:         t.State,
		Amount:        t.Amount,
		TransactionID: t.TransactionID,
		UserBalance:   newBalance,
	}, nil
}

func (s *TransactionService) CancelLatestOddTransactions() error {
	transactions, err := s.transactionRepo.GetLatestTransactions(19)
	if err != nil {
		return err
	}

	balance, err := s.balanceRepo.GetBalance()
	if err != nil {
		return err
	}

	var oddTransactionIDs []int
	var newBalance float64
	for i, transaction := range transactions {
		if i%2 != 0 {
			if transaction.State == stateWin {
				balance = balance - transaction.Amount
			} else {
				balance = balance + transaction.Amount
			}
			if balance >= 0 {
				oddTransactionIDs = append(oddTransactionIDs, transaction.ID)
				newBalance = balance
			} else {
				log.Printf("Balance is too small to cancel transaction. Aborting cancelling")
				break
			}
		}
	}

	if len(oddTransactionIDs) > 0 {
		if err := s.transactionRepo.CancelByIDs(oddTransactionIDs); err != nil {
			return err
		}

		if err := s.balanceRepo.SaveBalance(newBalance); err != nil {
			return err
		}

		log.Printf("%d transactions canceled, balance updated to %f", len(oddTransactionIDs), newBalance)
	} else {
		log.Printf("no transactions to cancel")
	}

	return nil
}

func calculateBalance(state string, balance float64, amount float64) float64 {
	var newBalance float64
	if state == stateWin {
		newBalance = balance + amount
	} else {
		newBalance = balance - amount
	}
	return newBalance
}
