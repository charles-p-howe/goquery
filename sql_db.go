package dataquery

import "reflect"

type SqlDb interface {
	Connection() interface{}
	Transaction() (Tx, error)
	Select(dest interface{}, stmt string, params ...interface{}) error
	Get(dest interface{}, stmt string, params ...interface{}) error
	Query(stmt string, params ...interface{}) (Rows, error)
	Insert(ds DataSet, rec interface{}, tx *Tx) error
}

type Rows interface {
	Columns() ([]string, error)
	ColumnTypes() ([]reflect.Type, error)
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}
