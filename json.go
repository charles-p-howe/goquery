package dataquery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/stoewer/go-strcase"
)

type jsonNullString struct {
	sql.NullString
}

func (v jsonNullString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.String)
}

type jsonNullInt32 struct {
	sql.NullInt32
}

func (v jsonNullInt32) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int32)
}

type jsonNullInt64 struct {
	sql.NullInt64
}

func (v jsonNullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

type jsonNullFloat64 struct {
	sql.NullFloat64
}

func (v jsonNullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

type jsonNullTime struct {
	sql.NullTime
	Fmt string
}

func (v jsonNullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	if v.Fmt != "" {
		//fmt.Println(v.Time.Format(v.Fmt))
		return json.Marshal(v.Time.Format(v.Fmt))
	}
	return json.Marshal(v.Time)
}

var jsonNullStringType = reflect.TypeOf(jsonNullString{})
var jsonNullInt32Type = reflect.TypeOf(jsonNullInt32{})
var jsonNullInt64Type = reflect.TypeOf(jsonNullInt64{})
var jsonNullFloat64Type = reflect.TypeOf(jsonNullFloat64{})
var jsonNullTimeType = reflect.TypeOf(jsonNullTime{})

var nullStringType = reflect.TypeOf(sql.NullString{})
var nullI32Type = reflect.TypeOf(sql.NullInt32{})
var nullI64Type = reflect.TypeOf(sql.NullInt64{})
var nullF64Type = reflect.TypeOf(sql.NullFloat64{})
var nullTimeType = reflect.TypeOf(sql.NullTime{})

var i32 int32
var i32Type = reflect.TypeOf(i32)
var i64 int64
var i64Type = reflect.TypeOf(i64)
var f64 float64
var f64Type = reflect.TypeOf(f64)
var str string
var strType = reflect.TypeOf(str)
var dte time.Time
var dateType = reflect.TypeOf(dte)

func RowsToJSON(rows *sql.Rows, toCamelCase bool, forceArray bool, dateFormat string, omitNull bool) ([]byte, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Column error: %v", err)
	}

	tt, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("Column type error: %v", err)
	}

	types := make([]reflect.Type, len(tt))
	for i, tp := range tt {
		st := tp.ScanType()
		if st == nil {
			return nil, fmt.Errorf("Scantype is null for column: %v", err)
		}
		switch st {
		case strType, nullStringType:
			types[i] = jsonNullStringType
		case i32Type, nullI32Type:
			types[i] = jsonNullInt32Type
		case i64Type, nullI64Type:
			types[i] = jsonNullInt64Type
		case f64Type, nullF64Type:
			types[i] = jsonNullFloat64Type
		case dateType, nullTimeType:
			types[i] = jsonNullTimeType
		default:
			types[i] = st

		}
	}

	values := make([]interface{}, len(tt))
	var builder strings.Builder
	count := 0

	for rows.Next() {
		if count > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('{')
		for i := range values {
			nv := reflect.New(types[i]) //.Interface()
			if types[i] == jsonNullTimeType {
				if dateFormat != "" {
					nv.Elem().Field(1).SetString(dateFormat)
				}
			}
			values[i] = nv.Interface()
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan values: %v", err)
		}
		for i, v := range values {
			jsonb, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			jsons := string(jsonb)
			if omitNull && jsons == "null" {
				continue
			}
			if i > 0 {
				builder.WriteRune(',')
			}
			fieldName := columns[i]
			if toCamelCase {
				fieldName = strcase.LowerCamelCase(fieldName)
			}
			builder.WriteString(fmt.Sprintf(`"%s":%s`, fieldName, jsons))
		}
		builder.WriteRune('}')
		count++
	}
	if count > 1 || forceArray {
		return []byte("[" + builder.String() + "]"), nil
	} else {
		return []byte(builder.String()), nil
	}
}
