package dataquery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type PgDataStore struct {
	DB                *pgx.Conn
	Config            *RdbmsConfig
	SequenceTemplate  SequenceTemplateFunction
	BindParamTemplate BindParamTemplateFunction
}

func NewPgxConnection(config *RdbmsConfig) (*pgx.Conn, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		config.Dbuser, config.Dbpass, config.Dbhost, config.Dbport, config.Dbname)
	con, err := pgx.Connect(context.Background(), dburl)
	return con, err
}
