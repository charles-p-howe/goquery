package dataquery

type RecordHandler func(interface{}) error

type BindParamTemplateFunction func(field string, i int) string
type SequenceTemplateFunction func(sequence string) string

type DataStore interface {
	Connection() interface{}
	GetSlice(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, panicOnErr bool) (interface{}, error)
	GetRecord(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, panicOnErr bool) (interface{}, error)
	GetJSON(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string) ([]byte, error)
	GetCSV(ds DataSet, key string, stmt string, suffix string, params []interface{}, statementAppends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string) (string, error)
	Select(ds DataSet) *FluentSelect
	//RecordsetIterator(s Select, handler RecordHandler)
	BeginTransaction() (Tx, error)
}
