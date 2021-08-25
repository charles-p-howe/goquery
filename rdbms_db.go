package goquery

type RdbmsDb interface {
	Connection() interface{}
	Transaction() (Tx, error)
	Select(dest interface{}, tx *Tx, stmt string, params ...interface{}) error
	Get(dest interface{}, tx *Tx, stmt string, params ...interface{}) error
	Query(tx *Tx, stmt string, params ...interface{}) (Rows, error)
	Insert(ds DataSet, rec interface{}, tx *Tx) error
	InsertStmt(ds DataSet) (string, error)
	Exec(tx *Tx, stmt string, params ...interface{}) error
	Batch() (Batch, error)
	SendBatch(batch Batch) BatchResult
}
