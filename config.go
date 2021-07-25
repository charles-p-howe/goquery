package dataquery

import (
	"os"
)

type RdbmsConfig struct {
	Dbuser string
	Dbpass string
	Dbhost string
	Dbport string
	Dbname string
}

func RdbmsConfigFromEnv() *RdbmsConfig {
	appConfig := new(RdbmsConfig)
	appConfig.Dbuser = os.Getenv("DBUSER")
	appConfig.Dbpass = os.Getenv("DBPASS")
	appConfig.Dbhost = os.Getenv("DBHOST")
	appConfig.Dbport = os.Getenv("DBPORT")
	appConfig.Dbname = os.Getenv("DBNAME")
	return appConfig
}
