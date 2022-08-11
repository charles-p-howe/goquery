package goquery

type FluentSelect struct {
	store       DataStore
	tx          *Tx
	qi          QueryInput
	dest        interface{}
	toCamelCase bool
	forceArray  bool
	dateFormat  string
	omitNull    bool
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
	return s
}

func (s *FluentSelect) CamelCase(useCamelCase bool) *FluentSelect {
	s.toCamelCase = useCamelCase
	return s
}

func (s *FluentSelect) DateFormat(dateFormat string) *FluentSelect {
	s.dateFormat = dateFormat
	return s
}

func (s *FluentSelect) OmitNull(omitnull bool) *FluentSelect {
	s.omitNull = omitnull
	return s
}

func (s *FluentSelect) ForceArray(forceArray bool) *FluentSelect {
	s.forceArray = forceArray
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

func (s *FluentSelect) Fetch() error {
	error := s.store.Fetch(s.tx, s.qi, s.dest)
	return error
}

func (s *FluentSelect) FetchRows() (Rows, error) {
	return s.store.FetchRows(s.tx, s.qi)
}

func (s *FluentSelect) FetchI() (interface{}, error) {
	dest := s.qi.DataSet.FieldSlice()
	error := s.store.Fetch(s.tx, s.qi, dest)
	return dest, error
}

func (s *FluentSelect) FetchJSON() ([]byte, error) {
	jsonOpts := JsonOpts{
		ToCamelCase: s.toCamelCase,
		OmitNull:    s.omitNull,
		ForceArray:  s.forceArray,
		DateFormat:  s.dateFormat,
	}
	return s.store.GetJSON(s.qi, jsonOpts)
}

func (s *FluentSelect) FetchCSV() (string, error) {
	csvOpts := CsvOpts{
		ToCamelCase: s.toCamelCase,
		DateFormat:  s.dateFormat,
		PrintHeader: true,
	}
	return s.store.GetCSV(s.qi, csvOpts)
}
