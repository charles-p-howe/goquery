package goquery

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
)

type TransactionFunction func(Tx)

var NoTx *Tx = nil

type Tx struct {
	tx interface{}
}

func (t Tx) PgxTx() *pgxpool.Tx {
	return t.tx.(*pgxpool.Tx)
}

func (t Tx) SqlXTx() *sqlx.Tx {
	return t.tx.(*sqlx.Tx)
}

func (t Tx) SqlTx() *sql.Tx {
	return t.tx.(*sql.Tx)
}

func (t Tx) Rollback() error {
	switch t.tx.(type) {
	case *sqlx.Tx:
		return t.tx.(*sqlx.Tx).Rollback()
	case *pgxpool.Tx:
		return t.tx.(*pgxpool.Tx).Rollback(context.Background())
	}
	return errors.New("invalid transaction type")
}

func (t Tx) Commit() error {
	switch t.tx.(type) {
	case *sqlx.Tx:
		return t.tx.(*sqlx.Tx).Commit()
	case *pgxpool.Tx:
		return t.tx.(*pgxpool.Tx).Commit(context.Background())
	}
	return errors.New("invalid transaction type")
}
