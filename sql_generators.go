package goquery

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func getSelectStatement(ds DataSet, key string, stmt string, suffix string, appends []interface{}) (string, error) {
	var ok bool
	if key != "" || stmt == "" {
		if ds == nil {
			return "", errors.New("Missing Dataset when referencing a statement key")
		}

		switch {
		case key != "":
			if stmt, ok = ds.Commands()[key]; !ok {
				return "", fmt.Errorf("unable to find statement for %s: %s", ds.Entity(), key)
			}
		case stmt == "":
			if stmt, ok = ds.Commands()[selectkey]; !ok {
				stmt = ToSelectStmt(ds)
				ds.PutCommand(selectkey, stmt)
			}
		}
	}

	stmt = fmt.Sprintf("%s %s", stmt, suffix)
	return fmt.Sprintf(stmt, appends...), nil
}

func ToSelectStmt(ds DataSet) string {
	fmt.Println("Building Statement")
	var fieldsBuilder strings.Builder
	fieldsBuilder.WriteString("select ")
	typ := reflect.TypeOf(ds.Attributes())
	fieldNum := typ.NumField()
	field := 0
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup("db"); ok && tagval != "_" {
			if field > 0 {
				fieldsBuilder.WriteRune(',')
			}
			fieldsBuilder.WriteString(tagval)
			field++
		}
	}
	fieldsBuilder.WriteString(fmt.Sprintf(" from %s", ds.Entity()))
	return fieldsBuilder.String()
}

func ToInsert(ds DataSet, dialect DbDialect) (string, error) {
	var fieldBuilder strings.Builder
	var bindBuilder strings.Builder
	typ := reflect.TypeOf(ds.Attributes())
	fieldNum := typ.NumField()
	fieldcount := 0
	paramcount := 0
	for i := 0; i < fieldNum; i++ {
		if dbfield, ok := typ.Field(i).Tag.Lookup("db"); ok && dbfield != "_" {
			if fieldcount > 0 {
				fieldBuilder.WriteRune(',')
				bindBuilder.WriteRune(',')
			}
			if idtype, ok := typ.Field(i).Tag.Lookup("dbid"); ok {
				if idtype != "AUTOINCREMENT" {
					if idsequence, ok := typ.Field(i).Tag.Lookup("idsequence"); ok {
						fieldBuilder.WriteString(dbfield)
						bindBuilder.WriteString(dialect.Seq(idsequence))
						fieldcount++
					} else {
						return "", errors.New("invalid id.  sequence type must have an 'idsequence' tag")
					}
				}
			} else {
				fieldBuilder.WriteString(dbfield)
				bindBuilder.WriteString(dialect.Bind(dbfield, paramcount))
				fieldcount++
				paramcount++
			}
		}
	}
	return fmt.Sprintf("insert into %s (%s) values (%s)", ds.Entity(), fieldBuilder.String(), bindBuilder.String()), nil
}

func ToUpdate(ds DataSet, bindTemplateFunction BindParamTemplateFunction) string {
	var fieldsBuilder strings.Builder
	var criteria string
	typ := reflect.TypeOf(ds.Attributes())
	fieldNum := typ.NumField()
	field := 0
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup("db"); ok && tagval != "_" {
			if field > 0 {
				fieldsBuilder.WriteRune(',')
			}
			if _, ok := typ.Field(i).Tag.Lookup("dbid"); ok {
				criteria = fmt.Sprintf(" where %s = %s", tagval, bindTemplateFunction(tagval, field))
				continue //skip id field
			}
			fieldsBuilder.WriteString(fmt.Sprintf("%s = %s", tagval, bindTemplateFunction(tagval, field)))
			field++
		}
	}
	return fmt.Sprintf("update %s set %s %s", ds.Entity(), fieldsBuilder.String(), criteria)
}

func IdField(ds DataSet) string {
	typ := reflect.TypeOf(ds.Attributes())
	fieldNum := typ.NumField()
	for i := 0; i < fieldNum; i++ {
		if tagval, ok := typ.Field(i).Tag.Lookup("db"); ok {
			if _, ok := typ.Field(i).Tag.Lookup("dbid"); ok {
				return tagval
			}
		}
	}
	return ""
}
