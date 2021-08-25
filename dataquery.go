package goquery

import (
	"fmt"
	"reflect"
)

const selectkey = "select"

//const updatekey = "update"
//const insertkey = "insert"

type Rows interface {
	Columns() ([]string, error)
	ColumnTypes() ([]reflect.Type, error)
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

type DataSet interface {
	Entity() string
	FieldSlice() interface{}
	Attributes() interface{}
	Commands() map[string]string
	PutCommand(key string, stmt string)
}

type TableDataSet struct {
	Name       string
	Schema     string //optional
	Statements map[string]string
	Fields     interface{}
}

func (t *TableDataSet) FieldSlice() interface{} {
	typ := reflect.TypeOf(t.Fields)
	slice := reflect.New(reflect.SliceOf(typ))
	return slice.Interface()
}

func (t *TableDataSet) Attributes() interface{} {
	return t.Fields
}

func (t *TableDataSet) Entity() string {
	if t.Schema != "" {
		return fmt.Sprintf("%s.%s", t.Schema, t.Name)
	}
	return t.Name
}

func (t *TableDataSet) Commands() map[string]string {
	return t.Statements
}

func (t *TableDataSet) PutCommand(key string, stmt string) {
	if t.Statements == nil {
		t.Statements = make(map[string]string)
	}
	t.Statements[selectkey] = stmt
}
