package dataquery

type FluentSelect struct {
	store       DataStore
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
	error := s.store.Fetch(s.qi, s.dest)
	return error
}

func (s *FluentSelect) FetchJSON() ([]byte, error) {
	s.qi.JsonOpts = &JsonOpts{
		ToCamelCase: s.toCamelCase,
		OmitNull:    s.omitNull,
		ForceArray:  s.forceArray,
		DateFormat:  s.dateFormat,
	}
	return s.store.GetJSON(s.qi)
}

func (s *FluentSelect) FetchCSV() (string, error) {
	s.qi.CsvOpts = &CsvOpts{
		ToCamelCase: s.toCamelCase,
		DateFormat:  s.dateFormat,
		PrintHeader: true,
	}
	return s.store.GetCSV(s.qi)
}
