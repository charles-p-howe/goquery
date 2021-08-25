package goquery

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx interface{}
}

func (t Tx) PgxTx() pgx.Tx {
	return t.tx.(pgx.Tx)
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
	case pgx.Tx:
		return t.tx.(pgx.Tx).Rollback(context.Background())
	}
	return errors.New("invalid transaction type")
}

func (t Tx) Commit() error {
	switch t.tx.(type) {
	case *sqlx.Tx:
		return t.tx.(*sqlx.Tx).Commit()
	case pgx.Tx:
		return t.tx.(pgx.Tx).Commit(context.Background())
	}
	return errors.New("invalid transaction type")
}

type TransactionFunction func(Tx)

/*
Transaction Wrapper.
DB Calls within the transaction should panic on error
*/

func Transaction(store DataStore, fn TransactionFunction) (err error) {
	var tx Tx
	tx, err = store.Transaction()
	if err != nil {
		log.Printf("Unable to start transaction: %s\n", err)
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			txerr := tx.Rollback()
			if txerr != nil {
				log.Printf("Unable to rollback from transaction: %s", err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Unable to commit transaction: %s", err)
			}
		}
	}()
	fn(tx)
	return err
}
