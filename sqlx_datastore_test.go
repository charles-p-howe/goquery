package dataquery

import (
	"reflect"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func setup(t *testing.T) DataStore {
	store := getStore(t)
	err := Transaction(store, func(tx Tx) {
		sqltx := tx.SqlXTx()
		sql := `create table fishing_spots(
			id serial not null primary key,
			location text)`
		sqltx.MustExec(sql)

		inserts := []string{
			`insert into fishing_spots (location) values ('Alpine Frove')`,
			`insert into fishing_spots (location) values ('Rivertown')`,
			`insert into fishing_spots (location) values ('Pine Island')`,
			`insert into fishing_spots (location) values (null)`,
		}
		for _, v := range inserts {
			sqltx.MustExec(v)
		}
	})
	if err != nil {
		t.Errorf("Setup error:%s\n", err)
	}
	return store
}

func teardown(store DataStore, t *testing.T) {
	err := Transaction(store, func(tx Tx) {
		sqltx := tx.SqlXTx()
		sqltx.MustExec("drop table fishing_spots")
	})
	if err != nil {
		t.Errorf("Failed to teardown test:%s\n", err)
	}
}

func getStore(t *testing.T) DataStore {
	config := RdbmsConfigFromEnv()
	db, err := NewSqlConnection(config)
	if err != nil {
		t.Errorf("Failed to connect to store:%s\n", err)
	}
	return &SqlDataStore{
		DB:                db,
		Config:            config,
		SequenceTemplate:  nil,
		BindParamTemplate: nil,
	}
}

func TestConnection(t *testing.T) {
	getStore(t)
}

func TestJson(t *testing.T) {
	correctResult := `[{"id":1,"location":"Alpine Frove"},{"id":2,"location":"Rivertown"},{"id":3,"location":"Pine Island"},{"id":4,"location":null}]`
	store := setup(t)
	defer teardown(store, t)

	json, err := store.Select(nil).
		Sql("select * from fishing_spots").
		OmitNull(false).
		FetchJSON()
	if err != nil {
		t.Errorf("Failed JSON Test: %s\n", err)
	}
	jsonstring := string(json)
	if jsonstring != correctResult {
		t.Errorf("Failed JSON Test: Got %s want %s", jsonstring, correctResult)
	}
}

func TestCsv(t *testing.T) {
	correctResult := `"id","location"
1,"Alpine Frove"
2,"Rivertown"
3,"Pine Island"
`
	store := setup(t)
	defer teardown(store, t)

	csv, err := store.Select(nil).
		Sql("select * from fishing_spots").
		FetchCSV()
	if err != nil {
		t.Errorf("Failed JSON Test: %s\n", err)
	}

	if csv != correctResult {
		t.Errorf("Failed CSV Test: Got %s want %s", csv, correctResult)
	}
}

type FishingSpot struct {
	ID       int32  `db:"id"`
	Location string `db:"location"`
}

func TestSlice(t *testing.T) {
	correctResult := []FishingSpot{
		{1, "Alpine Grove"},
		{2, "Rivertown"},
		{3, "Pine Island"},
	}

	store := setup(t)
	defer teardown(store, t)

	///////////autogenerate select///////////////////
	fsTbl := TableImpl{
		Name:   "fishing_spots",
		Fields: FishingSpot{},
	}

	res, err := store.Select(&fsTbl).FetchSlice()
	if err != nil {
		t.Errorf("Failed Slice Test:%s\n", err)
	}

	if reflect.DeepEqual(res, correctResult) {
		t.Errorf("Failed Slice Test: Got %v want %v", res, correctResult)
	}

	/////////////////add select statement///////////
	stmts := map[string]string{
		"named-select": `select * from fishing_spots`,
	}

	fsTbl.Statements = stmts
	res, err = store.Select(&fsTbl).
		StatementKey("named-select").
		FetchSlice()

	if err != nil {
		t.Errorf("Failed Slice Test:%s\n", err)
	}

	if reflect.DeepEqual(res, correctResult) {
		t.Errorf("Failed Slice Test: Got %v want %v", res, correctResult)
	}

}
