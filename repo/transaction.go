package repo

import (
	"context"
	"database/sql"

	"github.com/aidahputri/go-transaction/model"
)

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}

func (u *Transaction) Create(ctx context.Context, t model.Transaction) error {
	query := `INSERT INTO transaction (from_account, to_account, amount) VALUES ($1, $2, $3)`
	row := u.db.QueryRowContext(ctx, query, t.FromAccount, t.ToAccount, t.Amount)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}