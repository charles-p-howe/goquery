package goquery

type RdbmsDb interface {
	Connection() interface{}
	Transaction() (Tx, error)
	Select(dest interface{}, stmt string, params ...interface{}) error
	Get(dest interface{}, stmt string, params ...interface{}) error
	Query(stmt string, params ...interface{}) (Rows, error)
	Insert(ds DataSet, rec interface{}, tx *Tx) error
	InsertStmt(ds DataSet) (string, error)
	Exec(stmt string, params ...interface{}) error
	Batch() (Batch, error)
	SendBatch(batch Batch) BatchResult
}
