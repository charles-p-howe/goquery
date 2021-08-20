package dataquery

type RecordHandler func(interface{}) error

type BindParamTemplateFunction func(field string, i int) string
type SequenceTemplateFunction func(sequence string) string

type DbDialect struct {
	TableExistsStmt string
	Bind            BindParamTemplateFunction
	Seq             SequenceTemplateFunction
}

type DataStore interface {
	Connection() interface{}
	Transaction() (Tx, error)
	GetSlice(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, panicOnErr bool) (interface{}, error)
	GetRecord(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, panicOnErr bool) (interface{}, error)
	GetJSON(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string, omitNull bool) ([]byte, error)
	GetCSV(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string) (string, error)

	Select(ds DataSet) *FluentSelect
	Insert(ds DataSet) *FluentInsert
	//RecordsetIterator(s Select, handler RecordHandler)
	InsertRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	UpdateRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error
	Exec(stmt string, params ...interface{}) error
}
