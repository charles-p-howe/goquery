package dataquery

import (
	"fmt"
	"reflect"
)

const selectkey = "select"

//const updatekey = "update"
//const insertkey = "insert"

type DataSet interface {
	Entity() string
	FieldSlice() interface{}
	Attributes() interface{}
	Commands() map[string]string
	PutCommand(key string, stmt string)
}

type TableImpl struct {
	Name       string
	Schema     string //optional
	Statements map[string]string
	Fields     interface{}
}

func (t *TableImpl) FieldSlice() interface{} {
	typ := reflect.TypeOf(t.Fields)
	slice := reflect.New(reflect.SliceOf(typ))
	return slice.Interface()
}

func (t *TableImpl) Attributes() interface{} {
	return t.Fields
}

func (t *TableImpl) Entity() string {
	if t.Schema != "" {
		return fmt.Sprintf("%s.%s", t.Schema, t.Name)
	}
	return t.Name
}

func (t *TableImpl) Commands() map[string]string {
	return t.Statements
}

func (t *TableImpl) PutCommand(key string, stmt string) {
	if t.Statements == nil {
		t.Statements = make(map[string]string)
	}
	t.Statements[selectkey] = stmt
}
