package goquery

import "fmt"

var oracleDialect = DbDialect{
	TableExistsStmt: `SELECT count(*) FROM information_schema.tables WHERE  table_schema = $1 AND table_name = $2`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf(":%s", field)
	},
	Seq: func(sequence string) string {
		return fmt.Sprintf("nextval('%s')", sequence)
	},
	Url: func(config *RdbmsConfig) string {
		return fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s" libDir="%s"`,
			config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname, config.ExternalLib)
	},
}
