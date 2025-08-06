package goquery

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const defaultSSLMode = "disable"

// PoolMaxConnLifetime and PoolMaxConnIdle are string time duration representations
// as defined in ParseDuration in the stdlib time package
// the format consists of decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
type RdbmsConfig struct {
	Dbuser      string
	Dbpass      string
	Dbhost      string
	Dbport      string
	Dbname      string
	ExternalLib string
	OnInit      string
	DbDriver    string
	DbStore     string
	DbSSLMode   string

	PoolMaxConns        int
	PoolMinConns        int
	PoolMaxConnLifetime string //duration string
	PoolMaxConnIdle     string //duration string

	DbDriverSettings string
}

var sslModeMap = map[string]string{
	"disable":     "disable",
	"allow":       "allow",
	"prefer":      "prefer",
	"require":     "require",
	"verify-ca":   "verify-ca",
	"verify-full": "verify-full",
}

func RdbmsConfigFromEnv() *RdbmsConfig {
	dbConfig := new(RdbmsConfig)
	dbConfig.Dbuser = os.Getenv("DBUSER")
	dbConfig.Dbpass = os.Getenv("DBPASS")
	dbConfig.Dbhost = os.Getenv("DBHOST")
	dbConfig.Dbport = os.Getenv("DBPORT")
	dbConfig.Dbname = os.Getenv("DBNAME")
	dbConfig.DbDriver = os.Getenv("DBDRIVER")
	dbConfig.DbStore = os.Getenv("DBSTORE")
	dbConfig.ExternalLib = os.Getenv("EXTERNAL_LIB")
	dbConfig.DbSSLMode = os.Getenv("DBSSLMODE")

	if dbConfig.Dbport == "" {
		dbConfig.Dbport = "5432"
	}

	if dbConfig.DbSSLMode == "" {
		dbConfig.DbSSLMode = defaultSSLMode
	} else {
		if sslMode, ok := sslModeMap[strings.ToLower(dbConfig.DbSSLMode)]; ok {
			dbConfig.DbSSLMode = sslMode
		} else {
			log.Printf("Error parsing DBSSLMODE value of \"%s\":  Will fall back to default DBSSLMODE value.\n", dbConfig.DbSSLMode)
			dbConfig.DbSSLMode = defaultSSLMode
		}
	}

	maxConns := os.Getenv("POOLMAXCONNS")
	mc, err := strconv.Atoi(maxConns)
	if err != nil {
		log.Printf("Error parsing POOLMAXCONNS value of \"%s\":  Will fall back to default POOLMAXCONNS value.\n", maxConns)
	}
	dbConfig.PoolMaxConns = mc

	minConns := os.Getenv("POOLMINCONNS")
	mc, err = strconv.Atoi(minConns)
	if err != nil {
		log.Printf("Error parsing POOLMINCONNS value of \"%s\":  Will fall back to default POOLMINCONNS value.\n", minConns)
	}
	dbConfig.PoolMinConns = mc

	dbConfig.PoolMaxConnLifetime = os.Getenv("POOLMAXCONNLIFETIME")
	dbConfig.PoolMaxConnIdle = os.Getenv("POOLMAXCONNIDLE")

	return dbConfig
}

/*
MaxConnLifetime time.Duration

	// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	MaxConnLifetimeJitter time.Duration

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime time.Duration

	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	MaxConns int32

	// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// to create new connections.
	MinConns int32

*/

/*
func (db *DB) SetConnMaxIdleTime(d time.Duration)
func (db *DB) SetConnMaxLifetime(d time.Duration)
func (db *DB) SetMaxIdleConns(n int)
func (db *DB) SetMaxOpenConns(n int)
*/
