package goquery

import "fmt"

var pgDialect = DbDialect{
	TableExistsStmt: `SELECT count(*) FROM information_schema.tables WHERE  table_schema = $1 AND table_name = $2`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf("$%d", i+1)
	},
	Seq: func(sequence string) string {
		return fmt.Sprintf("nextval('%s')", sequence)
	},
}
