package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aidahputri/go-transaction/model"
)

type Account struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) *Account {
	return &Account{db: db}
}

func (u *Account) Create(ctx context.Context, a model.Account) error {
	query := `INSERT INTO accounts (name, account_number) VALUES ($1, $2)`
	row := u.db.QueryRowContext(ctx, query, a.Name, a.AccountNumber)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (u *Account) Get(ctx context.Context, accountNumber string) (model.Account, error) {
	query := `SELECT name, account_number, balance, blacklisted, under_watch FROM accounts WHERE account_number = $1`
	row := u.db.QueryRowContext(ctx, query, accountNumber)
	var acc model.Account

	err := row.Scan(&acc.Name, &acc.AccountNumber, &acc.Balance, &acc.Blacklisted, &acc.Underwatch)
	if err == sql.ErrNoRows {
		return model.Account{}, errors.New("account not found")
	}
	return acc, err
}

func (u *Account) Update(ctx context.Context, a model.Account) (model.Account, error) {
	query := `UPDATE accounts SET name = $1, balance = $3, blacklisted = $4, under_watch = $5 WHERE account_number = $6`
	res, err := u.db.ExecContext(ctx, query, a.Name, a.Balance, a.Blacklisted, a.Underwatch, a.AccountNumber)
	if err != nil {
		return model.Account{}, err
	}

	rowsAffected, err2 := res.RowsAffected()
	if err2 != nil {
		return model.Account{}, err2
	}
	if rowsAffected == 0 {
		return model.Account{}, sql.ErrNoRows
	}

	updatedAcc, err := u.Get(ctx, a.AccountNumber)
	if err != nil {
		return model.Account{}, err
	}

	return updatedAcc, nil
}