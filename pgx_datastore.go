package dataquery

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type PgxRows struct {
	rows pgx.Rows
}

func (p PgxRows) Columns() ([]string, error) {
	metadata := p.rows.FieldDescriptions()
	columns := make([]string, len(metadata))
	for _, f := range metadata {
		columns = append(columns, string(f.Name))
	}
	return columns, nil
}

func (p PgxRows) ColumnTypes() ([]reflect.Type, error) {
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

func (p PgxRows) Next() bool {
	return p.rows.Next()
}

func (p PgxRows) Scan(dest ...interface{}) error {
	return p.rows.Scan(dest)
}

func (p PgxRows) Close() error {
	p.rows.Close()
	return nil
}

type PgDataStore struct {
	DB                *pgx.Conn
	Config            *RdbmsConfig
	SequenceTemplate  SequenceTemplateFunction
	BindParamTemplate BindParamTemplateFunction
}

func NewPgxConnection(config *RdbmsConfig) (*pgx.Conn, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname)
	con, err := pgx.Connect(context.Background(), dburl)
	return con, err
}

func (pgs *PgDataStore) Connection() interface{} {
	return pgs.DB
}

func (pgs *PgDataStore) BeginTransaction() (Tx, error) {
	tx, err := pgs.DB.Begin(context.Background())
	return Tx{tx}, err
}

func (pgs *PgDataStore) GetSlice(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, panicOnErr bool) (interface{}, error) {
	ctx := context.Background()
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	data := ds.FieldSlice()
	if len(params) > 0 && params[0] != nil {
		err = pgxscan.Select(ctx, pgs.DB, data, sstmt, params...)
	} else {
		err = pgxscan.Select(ctx, pgs.DB, data, sstmt)
	}
	if err != nil && panicOnErr {
		panic(err)
	}
	return data, err
}

func (pgs *PgDataStore) GetRecord(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, panicOnErr bool) (interface{}, error) {
	ctx := context.Background()
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	typ := reflect.TypeOf(ds.Attributes())
	data := reflect.New(typ).Interface()
	if len(params) > 0 && params[0] != nil {
		err = pgxscan.Select(ctx, pgs.DB, data, sstmt, params...)
	} else {
		err = pgxscan.Select(ctx, pgs.DB, data, sstmt)
	}
	if err != nil && panicOnErr {
		panic(err)
	}
	return data, err
}

func (pgs *PgDataStore) GetJSON(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string, omitNull bool) ([]byte, error) {
	ctx := context.Background()
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return nil, err
	}
	//fmt.Println(sstmt)
	var rows pgx.Rows
	if len(params) > 0 && params[0] != nil {
		rows, err = pgs.DB.Query(ctx, sstmt, params...)
	} else {
		rows, err = pgs.DB.Query(ctx, sstmt)
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
	return RowsToJSON(PgxRows{rows}, toCamelCase, forceArray, dateFormat, omitNull)
}

func (pgs *PgDataStore) GetCSV(ds DataSet, key string, stmt string, suffix string, params []interface{}, appends []interface{}, toCamelCase bool, forceArray bool, panicOnErr bool, dateFormat string) (string, error) {
	ctx := context.Background()
	sstmt, err := getSelectStatement(ds, key, stmt, suffix, appends)
	if err != nil {
		return "", err
	}
	var rows pgx.Rows
	if len(params) > 0 && params[0] != nil {
		rows, err = pgs.DB.Query(ctx, sstmt, params...)
	} else {
		rows, err = pgs.DB.Query(ctx, sstmt)
	}
	if err != nil {
		log.Println(err)
		log.Println(sstmt)
		return "", err
	}
	defer rows.Close()
	return RowsToCSV(PgxRows{rows}, toCamelCase, dateFormat)
}

func (pgs *PgDataStore) Select(ds DataSet) *FluentSelect {
	s := FluentSelect{
		dataSet: ds,
		store:   pgs,
	}
	s.CamelCase(true)
	return &s
}
