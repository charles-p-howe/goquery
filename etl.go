package dataquery

type TransferOptions struct {
	CreateTable bool
	CommitSize  int
}

type ETL struct {
	source        DataStore
	sourceStmtKey string
	dest          DataStore
	destStmtKey   string
	options       TransferOptions
}

/*
func (etl *ETL) copyData(qi QueryInput) (err error) {
	rows, err := etl.source.FetchRows(qi)
	if err != nil {
		return err
	}
	defer rows.Close()
	typ := reflect.TypeOf(qi.DataSet.Attributes()).Elem()
	typeP := reflect.New(typ).Elem().Addr()
	structRef := typeP.Interface()
	var i int = 0


		for rows.Next() {
			i++
			err = rows.StructScan(structRef)
			if err != nil {
				panic(err)
			}
			if i%etl.options.CommitSize == 0 {
				err = etl.dest.Commit()
				if err != nil {
					panic(err)
				}
				err = etl.dest.StartTransaction()
				if err != nil {
					panic(err)
				}
			}
			etl.dest.CopyRow(table, i, structRef)
		}

	return nil
}
*/
