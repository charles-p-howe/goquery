package goquery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/stoewer/go-strcase"
)

var comma []byte = []byte(",")
var openobj []byte = []byte("{")
var closeobj []byte = []byte("}")
var openarray []byte = []byte("[")
var closearray []byte = []byte("]")

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

func RowsToJSON(builder io.Writer, rows Rows, toCamelCase bool, isArray bool, dateFormat string, omitNull bool) error {
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("Column error: %v", err)
	}

	tt, err := rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("Column type error: %v", err)
	}

	types := make([]reflect.Type, len(tt))
	for i, tp := range tt {
		//st := tp.ScanType()
		if tp == nil {
			return fmt.Errorf("Scantype is null for column: %v", err)
		}
		switch tp {
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
			types[i] = tp
		}
	}

	values := make([]interface{}, len(tt))
	//var builder strings.Builder
	count := 0

	if isArray {
		builder.Write(openarray)
	}
	for rows.Next() {
		if count > 0 {
			builder.Write(comma)
		}
		builder.Write(openobj)
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
			return fmt.Errorf("failed to scan values: %v", err)
		}
		for i, v := range values {
			jsonb, err := json.Marshal(v)
			if err != nil {
				return err
			}
			jsons := string(jsonb)
			if omitNull && jsons == "null" {
				continue
			}
			if i > 0 {
				builder.Write(comma)
			}
			fieldName := columns[i]
			if toCamelCase {
				fieldName = strcase.LowerCamelCase(fieldName)
			}
			builder.Write([]byte(fmt.Sprintf(`"%s":%s`, fieldName, jsons)))
		}
		builder.Write(closeobj)
		count++
	}

	if isArray {
		builder.Write(closearray)
	}
	return nil
}
