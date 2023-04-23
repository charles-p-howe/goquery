package goquery

import (
	"io"

	"github.com/jackc/pgconn"
)

type RecordHandler func(interface{}) error

type BindParamTemplateFunction func(field string, i int) string
type SequenceTemplateFunction func(sequence string) string
type UrlTemplateFunction func(config *RdbmsConfig) string

type DbDialect struct {
	TableExistsStmt string
	Bind            BindParamTemplateFunction
	Seq             SequenceTemplateFunction
	Url             UrlTemplateFunction
}

type QueryInput struct {
	DataSet      DataSet
	StatementKey string
	Statement    string
	Suffix       string
	BindParams   []interface{}
	StmtAppends  []interface{}
	PanicOnErr   bool
}

type InsertInput struct {
	Dataset    DataSet
	Records    interface{}
	Batch      bool
	BatchSize  int
	PanicOnErr bool
}

type JsonOpts struct {
	ToCamelCase bool
	IsArray     bool
	DateFormat  string
	OmitNull    bool
}

type CsvOpts struct {
	ToCamelCase bool
	DateFormat  string
	PrintHeader bool
}

type DataStore interface {
	Connection() interface{}
	NewTransaction() (Tx, error)
	Transaction(tf TransactionFunction) error
	Fetch(tx *Tx, input QueryInput, dest interface{}) error
	FetchRows(tx *Tx, input QueryInput) (Rows, error)
	GetJSON(writer io.Writer, input QueryInput, jo JsonOpts) error
	GetCSV(input QueryInput, co CsvOpts) (string, error)
	Select(stmt ...string) *FluentSelect
	Insert(ds DataSet) *FluentInsert
	//InsertRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	InsertRecs(tx *Tx, input InsertInput) error
	//UpdateRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	Exec(tx *Tx, stmt string, params ...interface{}) error
	MustExec(tx *Tx, stmt string, params ...interface{})
	//RecordsetIterator(s Select, handler RecordHandler)
}

type Batch interface {
	Queue(stmt string, params ...interface{})
}

type BatchResult interface {
	Exec() (pgconn.CommandTag, error)
	Close() error
}
