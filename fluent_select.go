package dataquery

type FluentSelect struct {
	store            DataStore
	dataSet          DataSet
	statementKey     string
	statementAppends []interface{}
	sql              string
	suffix           string
	params           []interface{}
	panicOnErr       bool
	//err              error
	toCamelCase bool
	forceArray  bool
	dateFormat  string
	omitNull    bool
}

func (s *FluentSelect) StatementKey(key string) *FluentSelect {
	s.statementKey = key
	return s
}

func (s *FluentSelect) Apply(vals ...interface{}) *FluentSelect {
	s.statementAppends = vals
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
	s.panicOnErr = panicOnErr
	return s
}

func (s *FluentSelect) Sql(stmt string) *FluentSelect {
	s.sql = stmt
	return s
}

func (s *FluentSelect) Suffix(suffix string) *FluentSelect {
	s.suffix = suffix
	return s
}

func (s *FluentSelect) Params(params ...interface{}) *FluentSelect {
	s.params = params
	return s
}

func (s *FluentSelect) FetchSlice() (interface{}, error) {
	recs, error := s.store.GetSlice(s.dataSet, s.statementKey, s.sql, s.suffix, s.params, s.statementAppends, s.panicOnErr)
	return recs, error
}

func (s *FluentSelect) FetchRow() (interface{}, error) {
	recs, error := s.store.GetRecord(s.dataSet, s.statementKey, s.sql, s.suffix, s.params, s.statementAppends, s.panicOnErr)
	return recs, error
}

func (s *FluentSelect) FetchJSON() ([]byte, error) {
	return s.store.GetJSON(s.dataSet, s.statementKey, s.sql, s.suffix, s.params, s.statementAppends, s.toCamelCase, s.forceArray, s.panicOnErr, s.dateFormat, s.omitNull)
}

func (s *FluentSelect) FetchCSV() (string, error) {
	return s.store.GetCSV(s.dataSet, s.statementKey, s.sql, s.suffix, s.params, s.statementAppends, s.toCamelCase, s.forceArray, s.panicOnErr, s.dateFormat)
}
