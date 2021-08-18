package dataquery

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type SqlRows struct {
	rows *sql.Rows
}

func (s SqlRows) Columns() ([]string, error) {
	return s.rows.Columns()
}

func (s SqlRows) ColumnTypes() ([]reflect.Type, error) {
	sts, err := s.rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	t := make([]reflect.Type, len(sts))
	for i := range sts {
		t[i] = sts[i].ScanType()
	}
	return t, nil
}

func (s SqlRows) Next() bool {
	return s.rows.Next()
}

func (s SqlRows) Scan(dest ...interface{}) error {
	return s.rows.Scan(dest...)
}

func (s SqlRows) Close() error {
	return s.rows.Close()
}

type SqlxDb struct {
	db *sqlx.DB
	//BeginTransaction BeginTransactionFunction
}

func NewSqlxConnection(config *RdbmsConfig) (SqlxDb, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname)
	con, err := sqlx.Connect("pgx", dburl)
	return SqlxDb{con}, err
}

func (sdb *SqlxDb) Connection() interface{} {
	return sdb.db
}

func (sdb *SqlxDb) Select(dest interface{}, stmt string, params ...interface{}) error {
	return sdb.db.Select(dest, stmt, params)
}

func (sdb *SqlxDb) Get(dest interface{}, stmt string, params ...interface{}) error {
	return sdb.db.Get(dest, stmt, params)
}

func (sdb *SqlxDb) Query(stmt string, params ...interface{}) (Rows, error) {
	rows, err := sdb.db.Query(stmt, params)
	return SqlRows{rows}, err
}

func (sdb *SqlxDb) BeginTransaction() (Tx, error) {
	tx, err := sdb.db.Beginx()
	return Tx{tx}, err
}
