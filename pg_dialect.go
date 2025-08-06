package goquery

import (
	"fmt"
	"log"
)

var pgDialect = DbDialect{
	TableExistsStmt: `SELECT count(*) FROM information_schema.tables WHERE  table_schema = $1 AND table_name = $2`,
	Bind: func(field string, i int) string {
		return fmt.Sprintf("$%d", i+1)
	},
	Seq: func(sequence string) string {
		return fmt.Sprintf("nextval('%s')", sequence)
	},
	Url: func(config *RdbmsConfig) string {
		if config.DbSSLMode == "" {
			config.DbSSLMode = defaultSSLMode
			log.Printf("No sslmode set, will fall back to default DBSSLMODE value %s. Set value in the dbconfig using DbSSLMode \n", defaultSSLMode)
		}
		return fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=%s",
			config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname, config.DbSSLMode)
	},
}
