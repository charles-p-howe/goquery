package goquery

import (
	"bufio"
	"bytes"
	"io"
)

type OutputFormat uint8

type FluentSelect struct {
	store       DataStore
	tx          *Tx
	qi          QueryInput
	qo          QueryOutput
	dest        interface{}
	rowFunction RowFunction
}

func (s *FluentSelect) DataSet(ds DataSet) *FluentSelect {
	s.qi.DataSet = ds
	return s
}

func (s *FluentSelect) Tx(tx *Tx) *FluentSelect {
	s.tx = tx
	return s
}

func (s *FluentSelect) StatementKey(key string) *FluentSelect {
	s.qi.StatementKey = key
	return s
}

func (s *FluentSelect) Apply(vals ...interface{}) *FluentSelect {
	s.qi.StmtAppends = vals
	return s
}

func (s *FluentSelect) Dest(dest interface{}) *FluentSelect {
	s.dest = dest
	s.qo.OutputFormat = DEST
	return s
}

func (s *FluentSelect) CamelCase(useCamelCase bool) *FluentSelect {
	s.qo.Options.ToCamelCase = useCamelCase
	return s
}

func (s *FluentSelect) DateFormat(dateFormat string) *FluentSelect {
	s.qo.Options.DateFormat = dateFormat
	return s
}

func (s *FluentSelect) OmitNull(omitnull bool) *FluentSelect {
	s.qo.Options.OmitNull = omitnull
	return s
}

func (s *FluentSelect) IsJsonArray(isJsonArray bool) *FluentSelect {
	s.qo.Options.IsArray = isJsonArray
	return s
}

func (s *FluentSelect) PanicOnErr(panicOnErr bool) *FluentSelect {
	s.qi.PanicOnErr = panicOnErr
	return s
}

func (s *FluentSelect) Suffix(suffix string) *FluentSelect {
	s.qi.Suffix = suffix
	return s
}

func (s *FluentSelect) Params(params ...interface{}) *FluentSelect {
	s.qi.BindParams = params
	return s
}

func (s *FluentSelect) OutputJson(writer io.Writer) *FluentSelect {
	s.qo.Writer = writer
	s.qo.OutputFormat = JSON
	return s
}

func (s *FluentSelect) OutputCsv(writer io.Writer) *FluentSelect {
	s.qo.Writer = writer
	s.qo.OutputFormat = CSV
	return s
}

func (s *FluentSelect) ForEachRow(rf RowFunction) *FluentSelect {
	s.rowFunction = rf
	return s
}

func (s *FluentSelect) Fetch() error {
	error := s.store.Fetch(s.tx, s.qi, s.qo, s.dest)
	return error
}

func (s *FluentSelect) FetchRows() (Rows, error) {
	return s.store.FetchRows(s.tx, s.qi)
}

// Deprecated: This method will be removed in the next version.  Use Fetch()
func (s *FluentSelect) FetchI() (interface{}, error) {
	dest := s.qi.DataSet.FieldSlice()
	error := s.store.Fetch(s.tx, s.qi, s.qo, dest)
	return dest, error
}

// Deprecated: This method will be removed in the next version.  Use Fetch()
func (s *FluentSelect) FetchJSON() ([]byte, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err := s.store.GetJSON(writer, s.qi, s.qo.Options)
	writer.Flush()
	return b.Bytes(), err
}

// Deprecated: This method will be removed in the next version.  Use Fetch()
/*
func (s *FluentSelect) FetchCSV() (string, error) {

	return s.store.GetCSV(s.qi, s.qo.Options)
}
*/
