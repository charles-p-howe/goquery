package dataquery

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

//@TODO panic on error is not complete
//implements the datastore interface

func NewRdbmsDataStore(config *RdbmsConfig) (DataStore, error) {
	switch config.DbStore {
	case "pgx":
		db, err := NewPgxConnection(config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to connect to pgx datastore: %s", err))
		}
		return &RdbmsDataStore{&db}, nil
	case "sqlx":
		db, err := NewSqlxConnection(config)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to connect to pgx datastore: %s", err))
		}
		return &RdbmsDataStore{&db}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported store type: %s", config.DbStore))
	}
}

type RdbmsDataStore struct {
	db RdbmsDb
}

func (sds *RdbmsDataStore) RdbmsDb() RdbmsDb {
	return sds.db
}

func (sds *RdbmsDataStore) Connection() interface{} {
	return sds.db.Connection()
}

func (sds *RdbmsDataStore) Transaction() (Tx, error) {
	return sds.db.Transaction()
}

func (sds *RdbmsDataStore) Fetch(qi QueryInput, dest interface{}) error {
	sstmt, err := getSelectStatement(qi.DataSet, qi.StatementKey, qi.Statement, qi.Suffix, qi.StmtAppends)
	if err != nil {
		return err
	}
	if len(qi.BindParams) > 0 && qi.BindParams[0] != nil {
		err = sds.db.Select(dest, sstmt, qi.BindParams...)
	} else {
		err = sds.db.Select(dest, sstmt)
	}
	if err != nil && qi.PanicOnErr {
		panic(err)
	}
	return err
}

func (sds *RdbmsDataStore) FetchRows(qi QueryInput) (Rows, error) {
	sstmt, err := getSelectStatement(qi.DataSet, qi.StatementKey, qi.Statement, qi.Suffix, qi.StmtAppends)
	if err != nil {
		return nil, err
	}

	if len(qi.BindParams) > 0 && qi.BindParams[0] != nil {
		return sds.db.Query(sstmt, qi.BindParams...)
	} else {
		return sds.db.Query(sstmt)
	}

}

func (sds *RdbmsDataStore) GetJSON(qi QueryInput, jo JsonOpts) ([]byte, error) {
	rows, err := sds.FetchRows(qi)
	if err != nil {
		log.Println(err)
		if qi.PanicOnErr {
			panic(err)
		}
		return nil, err
	}
	defer rows.Close()
	return RowsToJSON(rows, jo.ToCamelCase, jo.ForceArray, jo.DateFormat, jo.OmitNull)
}

func (sds *RdbmsDataStore) GetCSV(qi QueryInput, co CsvOpts) (string, error) {
	rows, err := sds.FetchRows(qi)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer rows.Close()
	return RowsToCSV(rows, co.ToCamelCase, co.DateFormat)
}

func (sds *RdbmsDataStore) InsertRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error {
	rval := reflect.ValueOf(recs)
	rrecs := reflect.Indirect(rval)
	if rrecs.Kind() == reflect.Slice {
		if batch {
			sds.insertBatch(ds, rrecs, batchSize)
		} else {
			if tx == nil {
				return sds.insertNewTrans(ds, rrecs)
			} else {
				return sds.insert(ds, rrecs, tx)
			}
		}
	} else {
		return sds.db.Insert(ds, recs, tx)
	}
	return nil
}

func (sds *RdbmsDataStore) Exec(stmt string, params ...interface{}) error {
	return sds.db.Exec(stmt, params)
}

func (sds *RdbmsDataStore) insertNewTrans(ds DataSet, rrecs reflect.Value) error {
	err := Transaction(sds, func(tx Tx) {
		err := sds.insert(ds, rrecs, &tx)
		if err != nil {
			panic(err)
		}
	})
	return err
}

func (sds *RdbmsDataStore) insert(ds DataSet, rrecs reflect.Value, tx *Tx) error {
	for i := 0; i < rrecs.Len(); i++ {
		err := sds.db.Insert(ds, rrecs.Index(i).Interface(), tx)
		if err != nil {
			log.Printf("Failed to insert: %s\n", err)
			return err
		}
	}
	return nil
}

func (sds *RdbmsDataStore) insertBatch(ds DataSet, rrecs reflect.Value, batchSize int) error {
	batch, err := sds.db.Batch()
	if err != nil {
		return err
	}

	stmt, err := sds.db.InsertStmt(ds)
	if err != nil {
		return err
	}

	for i := 0; i < rrecs.Len(); i++ {
		rec := rrecs.Index(i).Interface()
		batch.Queue(stmt, StructToIArray(rec)...)
		if i >= batchSize {
			sds.db.SendBatch(batch)
			batch, err = sds.db.Batch()
			if err != nil {
				return err
			}
		}
	}
	sds.db.SendBatch(batch)
	return nil
}

/*
func hasErr(br BatchResult) error {
	for i := 0; i < 10; i++ {
		ct, err := br
	}
}
*/

/*
func (sds *RdbmsDataStore) InsertRecs(ds DataSet, recs interface{}, batch bool, batchSize int, tx *Tx) error {
	rval := reflect.ValueOf(recs)
	rrecs := reflect.Indirect(rval)
	if rval.Kind() == reflect.Slice {
		if tx == nil {
			err := Transaction(sds, func(tx Tx) {
				for i := 0; i < rrecs.Len(); i++ {
					err := sds.db.Insert(ds, rrecs.Index(i).Interface(), &tx)
					if err != nil {
						panic(err)
					}
				}
			})
			return err
		} else {
			for i := 0; i < rrecs.Len(); i++ {
				err := sds.db.Insert(ds, rrecs.Index(i).Interface(), tx)
				if err != nil {
					log.Printf("Failed to insert: %s\n", err)
					return err
				}
			}
		}
	} else {
		return sds.db.Insert(ds, recs, tx)
	}
	return nil
}
*/

/*
func (sds *SqlDataStore) Insert(ds DataSet, val interface{}, retval interface{}, tx *sqlx.Tx) error {
	var err error
	if retval != nil {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, func(field string, i int) string { return fmt.Sprintf("$%d", i) })
		if err != nil {
			return err
		}
		stmt = stmt + " returning id"
		fmt.Println(stmt)
		if tx == nil {
			err = sds.DB.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		} else {
			err = tx.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		}
		if err != nil {
			return err
		}
	} else {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, sds.BindParamTemplate)
		if err != nil {
			return err
		}
		_, err = sds.DB.NamedExec(stmt, val)
		fmt.Println(err)
	}
	//@TODO this error is getting shadowed by the inner error...need to fix
	return err
}
*/

/*
func (sds *SqlDataStore) insertRec(rec interface{}, retval interface{}) error {
	var err error
	if retval != nil {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, func(field string, i int) string { return fmt.Sprintf("$%d", i) })
		if err != nil {
			return err
		}
		stmt = stmt + " returning id"
		fmt.Println(stmt)
		if tx == nil {
			err = sds.DB.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		} else {
			err = tx.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		}
		if err != nil {
			return err
		}
	} else {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, sds.BindParamTemplate)
		if err != nil {
			return err
		}
		_, err = sds.DB.NamedExec(stmt, val)
		fmt.Println(err)
	}
	//@TODO this error is getting shadowed by the inner error...need to fix
	return err
}
*/

func (sds *RdbmsDataStore) Select(stmt ...string) *FluentSelect {
	stmts := ""
	if len(stmt) > 0 && stmt[0] != "" {
		stmts = stmt[0]
	}
	s := FluentSelect{
		qi: QueryInput{
			Statement: stmts,
		},
		store: sds,
	}
	s.CamelCase(true)
	return &s
}

/*
func (sds *SqlDataStore) Select(ds DataSet) *FluentSelect {
	s := FluentSelect{
		qi: QueryInput{
			DataSet: ds,
		},
		store: sds,
	}
	s.CamelCase(true)
	return &s
}
*/

func (sds *RdbmsDataStore) Insert(ds DataSet) *FluentInsert {
	fi := FluentInsert{
		ds:    ds,
		store: sds,
	}
	return &fi
}

/*
type SqlxDataStore struct {
	DB                *sqlx.DB
	Config            *RdbmsConfig
	SequenceTemplate  SequenceTemplateFunction
	BindParamTemplate BindParamTemplateFunction
}

func NewSqlxConnection(config *RdbmsConfig) (*sqlx.DB, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname)
	con, err := sqlx.Connect("pgx", dburl)
	return con, err
}

func (sds *SqlxDataStore) Connection() interface{} {
	return sds.DB
}

func (sds *SqlxDataStore) BeginTransaction() (Tx, error) {
	tx, err := sds.DB.Beginx()
	return Tx{tx}, err
}

func (sds *SqlxDataStore) GetSlice(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, panicOnErr bool) (interface{}, error) {
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	data := ds.FieldSlice()
	if len(params) > 0 && params[0] != nil {
		err = sds.DB.Select(data, sstmt, params...)
	} else {
		err = sds.DB.Select(data, sstmt)
	}
	if err != nil && panicOnErr {
		panic(err)
	}
	return data, err
}

func (sds *SqlxDataStore) GetRecord(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, panicOnErr bool) (interface{}, error) {
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	typ := reflect.TypeOf(ds.Attributes())
	data := reflect.New(typ).Interface()
	if len(params) > 0 && params[0] != nil {
		err = sds.DB.Get(data, sstmt, params...)
	} else {
		err = sds.DB.Get(data, sstmt)
	}
	if err != nil && panicOnErr {
		panic(err)
	}
	return data, err
}

func (sds *SqlxDataStore) GetJSON(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string, omitNull bool) ([]byte, error) {
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	//fmt.Println(sstmt)
	var rows *sql.Rows
	if len(params) > 0 && params[0] != nil {
		rows, err = sds.DB.Query(sstmt, params...)
	} else {
		rows, err = sds.DB.Query(sstmt)
	}
	if err != nil {
		log.Println(err)
		log.Println(sstmt)
		if panicOnErr {
			panic(err)
		}
		return nil, err
	}
	defer rows.Close()
	return RowsToJSON(SqlRows{rows}, toCamelCase, forceArray, dateFormat, omitNull)
}

func (sds *SqlxDataStore) GetCSV(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string) (string, error) {
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return "", err
	}
	var rows *sql.Rows
	if len(params) > 0 && params[0] != nil {
		rows, err = sds.DB.Query(sstmt, params...)
	} else {
		rows, err = sds.DB.Query(sstmt)
	}
	if err != nil {
		log.Println(err)
		log.Println(sstmt)
		return "", err
	}
	defer rows.Close()
	return RowsToCSV(SqlRows{rows}, toCamelCase, dateFormat)
}

func (sds *SqlxDataStore) Select(ds DataSet) *FluentSelect {
	s := FluentSelect{
		dataSet: ds,
		store:   sds,
	}
	s.CamelCase(true)
	return &s
}

func (sds *SqlxDataStore) Insert(ds DataSet, val interface{}, retval interface{}, tx *sqlx.Tx) error {
	var err error
	if retval != nil {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, func(field string, i int) string { return fmt.Sprintf("$%d", i) })
		if err != nil {
			return err
		}
		stmt = stmt + " returning id"
		fmt.Println(stmt)
		if tx == nil {
			err = sds.DB.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		} else {
			err = tx.Get(retval, stmt, ValsAsInterfaceArray2(val, []string{"ID"}, "db", []string{"_"})...)
		}
		if err != nil {
			return err
		}
	} else {
		stmt, err := ToInsert(ds, sds.SequenceTemplate, sds.BindParamTemplate)
		if err != nil {
			return err
		}
		fmt.Println(stmt)
		_, err = sds.DB.NamedExec(stmt, val)
		if err != nil {
			return err
		}
	}
	//@TODO this error is getting shadowed by the inner error...need to fix
	return err
}

func (sds *SqlxDataStore) Update(ds DataSet, val interface{}) error {
	stmt := ToUpdate(ds, sds.BindParamTemplate)
	fmt.Println(stmt)
	_, err := sds.DB.NamedExec(stmt, val)
	return err
}

func (sds *SqlxDataStore) Delete(ds DataSet, id interface{}) error {
	template := "delete from %s where %s = $1"
	idfield := IdField(ds)
	stmt := fmt.Sprintf(template, ds.Entity(), idfield)
	_, err := sds.DB.Exec(stmt, id)
	return err
}
*/
