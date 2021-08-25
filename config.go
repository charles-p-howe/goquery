package goquery

import (
	"os"
)

type RdbmsConfig struct {
	Dbuser      string
	Dbpass      string
	Dbhost      string
	Dbport      string
	Dbname      string
	ExternalLib string
	DbDriver    string
	DbStore     string
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
	return dbConfig
}
