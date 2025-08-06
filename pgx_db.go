package goquery

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxExecResult struct {
	res pgconn.CommandTag
}

func (per PgxExecResult) RowsAffected() int64 {
	return per.res.RowsAffected()
}

// @TODO figure out how to handle command tags and then we only need a single Execr interface
type PgxExecr interface {
	Exec(ctx context.Context, stmt string, params ...interface{}) (pgconn.CommandTag, error)
}

type PgxRows struct {
	rows       pgx.Rows
	rowScanner *pgxscan.RowScanner
}

func (p *PgxRows) Columns() ([]string, error) {
	metadata := p.rows.FieldDescriptions()
	columns := make([]string, len(metadata))
	for i, f := range metadata {
		columns[i] = string(f.Name)
	}
	return columns, nil
}

func (p *PgxRows) ColumnTypes() ([]reflect.Type, error) {
	metadata := p.rows.FieldDescriptions()
	t := make([]reflect.Type, len(metadata))
	for i, fd := range metadata {
		switch fd.DataTypeOID {
		case pgtype.Float8OID:
			t[i] = reflect.TypeOf(float64(0))
		case pgtype.Float4OID:
			t[i] = reflect.TypeOf(float32(0))
		case pgtype.Int8OID:
			t[i] = reflect.TypeOf(int64(0))
		case pgtype.Int4OID:
			t[i] = reflect.TypeOf(int32(0))
		case pgtype.Int2OID:
			t[i] = reflect.TypeOf(int16(0))
		case pgtype.BoolOID:
			t[i] = reflect.TypeOf(false)
		case pgtype.NumericOID:
			t[i] = reflect.TypeOf(float64(0))
		case pgtype.DateOID, pgtype.TimestampOID, pgtype.TimestamptzOID:
			t[i] = reflect.TypeOf(time.Time{})
		case pgtype.ByteaOID:
			t[i] = reflect.TypeOf([]byte(nil))
		default:
			t[i] = reflect.TypeOf("")
		}
	}
	return t, nil
}

func (p *PgxRows) Next() bool {
	return p.rows.Next()
}

func (p *PgxRows) Scan(dest ...interface{}) error {
	return p.rows.Scan(dest...)
}

func (p *PgxRows) ScanStruct(dest interface{}) error {
	if p.rowScanner == nil {
		p.rowScanner = pgxscan.NewRowScanner(p.rows)
	}
	return p.rowScanner.Scan(dest)
}

func (p *PgxRows) Close() error {
	p.rows.Close()
	return nil
}

/*
type PgxBatch struct {
	BatchSize int
	batch     *pgx.Batch
}

func(pb *PgxBatch) Queue(stmt string, params ...interface{}){
	pb.batch.Queue(stmt,params)
}
*/

type PgxDb struct {
	db      *pgxpool.Pool
	dialect DbDialect
}

func NewPgxConnection(config *RdbmsConfig) (PgxDb, error) {
	if config.DbSSLMode == "" {
		config.DbSSLMode = defaultSSLMode
		log.Printf("No sslmode set, will fall back to default DBSSLMODE value %s. Set value in the dbconfig using DbSSLMode \n", defaultSSLMode)
	}
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=%s",
		config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname, config.DbSSLMode)

	if config.PoolMaxConns > 0 {
		dburl = fmt.Sprintf("%s %s=%d", dburl, "pool_max_conns", config.PoolMaxConns)
	}
	if config.PoolMinConns > 0 {
		dburl = fmt.Sprintf("%s %s=%d", dburl, "pool_min_conns", config.PoolMinConns)
	}
	if config.PoolMaxConnLifetime != "" {
		dburl = fmt.Sprintf("%s %s=%s", dburl, "pool_max_conn_lifetime", config.PoolMaxConnLifetime)
	}
	if config.PoolMaxConnIdle != "" {
		dburl = fmt.Sprintf("%s %s=%s", dburl, "pool_max_conn_idle_time", config.PoolMaxConnIdle)
	}

	con, err := pgxpool.Connect(context.Background(), dburl)
	return PgxDb{con, pgDialect}, err
}

func (pdb *PgxDb) Connection() interface{} {
	return pdb.db
}

func (pdb *PgxDb) querier(tx *Tx) pgxscan.Querier {
	if tx != nil {
		return tx.PgxTx()
	}
	return pdb.db
}

func (pdb *PgxDb) execr(tx *Tx) PgxExecr {
	if tx != nil {
		return tx.PgxTx()
	}
	return pdb.db
}

func (pdb *PgxDb) Select(dest interface{}, tx *Tx, stmt string, params ...interface{}) error {
	return pgxscan.Select(context.Background(), pdb.querier(tx), dest, stmt, params...)
}

func (pdb *PgxDb) Get(dest interface{}, tx *Tx, stmt string, params ...interface{}) error {
	return pgxscan.Get(context.Background(), pdb.querier(tx), dest, stmt, params...)
}

func (pdb *PgxDb) Query(tx *Tx, stmt string, params ...interface{}) (Rows, error) {
	rows, err := pdb.querier(tx).Query(context.Background(), stmt, params...)
	return &PgxRows{rows, nil}, err
}

// @DEPRICATED
func (pdb *PgxDb) Exec(tx *Tx, stmt string, params ...interface{}) error {
	_, err := pdb.execr(tx).Exec(context.Background(), stmt, params...)
	return err
}

func (pdb *PgxDb) Execr(tx *Tx, stmt string, params ...interface{}) (ExecResult, error) {
	ct, err := pdb.execr(tx).Exec(context.Background(), stmt, params...)
	return PgxExecResult{ct}, err
}

func (pdb *PgxDb) MustExec(tx *Tx, stmt string, params ...interface{}) {
	_, err := pdb.execr(tx).Exec(context.Background(), stmt, params...)
	if err != nil {
		panic(err)
	}
}

func (pdb *PgxDb) MustExecr(tx *Tx, stmt string, params ...interface{}) ExecResult {
	ct, err := pdb.execr(tx).Exec(context.Background(), stmt, params...)
	if err != nil {
		panic(err)
	}
	return PgxExecResult{ct}
}

func (pdb *PgxDb) Batch() (Batch, error) {
	return &pgx.Batch{}, nil
}

func (pdb *PgxDb) SendBatch(batch Batch) BatchResult {
	pb := batch.(*pgx.Batch)
	br := pdb.db.SendBatch(context.Background(), pb)
	br.Close()
	return br
}

func (pdb *PgxDb) InsertStmt(ds DataSet) (string, error) {
	return ToInsert(ds, pdb.dialect)
}

func (pdb *PgxDb) Insert(ds DataSet, rec interface{}, tx *Tx) error {

	var err error
	var stmt string
	var ok bool

	if stmt, ok = ds.Commands()["insert"]; !ok {
		stmt, err = ToInsert(ds, pdb.dialect)
	}
	if err != nil {
		return err
	}
	params := StructToIArray(rec)
	if tx == nil {
		_, err = pdb.db.Exec(context.Background(), stmt, params...)
		return err
	} else {
		pgxtx := tx.PgxTx()
		_, err = pgxtx.Exec(context.Background(), stmt, params...)
		return err
	}

}

func (pdb *PgxDb) Transaction() (Tx, error) {
	tx, err := pdb.db.Begin(context.Background())
	return Tx{tx}, err
}
