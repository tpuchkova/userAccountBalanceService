package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

const userID = 1

func (r *BalancePostgres) SaveBalance(balance float64) error {
	query := fmt.Sprintf("UPDATE %s SET balance = '%f' WHERE id = '%d'", userBalanceTable, balance, userID)
	_, err := r.db.Exec(query)
	return err
}

func (r *BalancePostgres) GetBalance() (float64, error) {
	var balance float64
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = '%d'", userBalanceTable, userID)
	if err := r.db.QueryRow(query).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
