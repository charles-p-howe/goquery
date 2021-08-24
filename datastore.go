package dataquery

import "github.com/jackc/pgconn"

type RecordHandler func(interface{}) error

type BindParamTemplateFunction func(field string, i int) string
type SequenceTemplateFunction func(sequence string) string

type DbDialect struct {
	TableExistsStmt string
	Bind            BindParamTemplateFunction
	Seq             SequenceTemplateFunction
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

type JsonOpts struct {
	ToCamelCase bool
	ForceArray  bool
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
	Transaction() (Tx, error)
	Fetch(input QueryInput, dest interface{}) error
	FetchRows(input QueryInput) (Rows, error)
	GetJSON(input QueryInput, jo JsonOpts) ([]byte, error)
	GetCSV(input QueryInput, co CsvOpts) (string, error)
	Select(stmt ...string) *FluentSelect
	Insert(ds DataSet) *FluentInsert
	InsertRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	//UpdateRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	Exec(stmt string, params ...interface{}) error
	//RecordsetIterator(s Select, handler RecordHandler)
}

type Batch interface {
	Queue(stmt string, params ...interface{})
}

type BatchResult interface {
	Exec() (pgconn.CommandTag, error)
	Close() error
}
