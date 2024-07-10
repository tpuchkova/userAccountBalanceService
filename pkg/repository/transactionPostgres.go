package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"userAccountBalanceService/pkg/repository/model"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) GetLatestTransactions(count int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	query := fmt.Sprintf("SELECT id, state, amount FROM %s WHERE canceled = false ORDER BY created_at DESC LIMIT %d", transactionsTable, count)
	if err := r.db.Select(&transactions, query); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionPostgres) CancelByIDs(idsToCancel []int) error {
	idStrings := make([]string, len(idsToCancel))
	for i, id := range idsToCancel {
		idStrings[i] = fmt.Sprintf("%d", id)
	}

	ids := strings.Join(idStrings, ", ")

	query := fmt.Sprintf("UPDATE %s SET canceled = true WHERE id IN (%s)", transactionsTable, ids)
	_, err := r.db.Exec(query)

	return err
}

func (r *TransactionPostgres) AddTransaction(transaction model.Transaction) (*model.Transaction, error) {
	var id int
	var state string
	var amount float64
	var transactionID string
	createQuery := fmt.Sprintf("INSERT INTO %s (state, amount, transaction_id) values ($1, $2, $3) RETURNING id, state, amount, transaction_id", transactionsTable)

	row := r.db.QueryRow(createQuery, transaction.State, transaction.Amount, transaction.TransactionID)
	if err := row.Scan(&id, &state, &amount, &transactionID); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:            id,
		State:         state,
		Amount:        amount,
		TransactionID: transactionID,
	}, nil
}
