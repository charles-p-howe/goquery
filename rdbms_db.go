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
	Execr(tx *Tx, stmt string, params ...interface{}) (ExecResult, error)
	MustExec(tx *Tx, stmt string, params ...interface{})
	MustExecr(tx *Tx, stmt string, params ...interface{}) ExecResult
	Batch() (Batch, error)
	SendBatch(batch Batch) BatchResult
}
