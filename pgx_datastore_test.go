package dataquery

import (
	"context"
	"reflect"
	"testing"
)

func pgxsetup(t *testing.T) DataStore {
	ctx := context.Background()
	store := getPgxStore(t)
	err := Transaction(store, func(tx Tx) {
		pgxtx := tx.PgxTx()
		sql := `create table fishing_spots(
			id serial not null primary key,
			location text)`
		_, err := pgxtx.Exec(ctx, sql)
		if err != nil {
			panic(err)
		}

		inserts := []string{
			`insert into fishing_spots (location) values ('Alpine Frove')`,
			`insert into fishing_spots (location) values ('Rivertown')`,
			`insert into fishing_spots (location) values ('Pine Island')`,
			`insert into fishing_spots (location) values (null)`,
		}
		for _, v := range inserts {
			_, err := pgxtx.Exec(ctx, v)
			if err != nil {
				panic(err)
			}
		}
	})
	if err != nil {
		t.Errorf("Setup error:%s\n", err)
	}
	return store
}

func pgxteardown(store DataStore, t *testing.T) {
	ctx := context.Background()
	err := Transaction(store, func(tx Tx) {
		pgxtx := tx.PgxTx()
		_, err := pgxtx.Exec(ctx, "drop table fishing_spots")
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		t.Errorf("Failed to teardown test:%s\n", err)
	}
}

func getPgxStore(t *testing.T) DataStore {
	config := RdbmsConfigFromEnv()
	db, err := NewPgxConnection(config)
	if err != nil {
		t.Errorf("Failed to connect to store:%s\n", err)
	}
	return &PgDataStore{
		DB:                db,
		Config:            config,
		SequenceTemplate:  nil,
		BindParamTemplate: nil,
	}
}

func TestPgxConnection(t *testing.T) {
	getPgxStore(t)
}

func TestPgxJson(t *testing.T) {
	correctResult := `[{"id":1,"location":"Alpine Frove"},{"id":2,"location":"Rivertown"},{"id":3,"location":"Pine Island"},{"id":4,"location":null}]`
	store := pgxsetup(t)
	defer pgxteardown(store, t)

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

func TestPgxSlice(t *testing.T) {
	ap := "Alpine Frove"
	rt := "Rivertown"
	pi := "Pine Island"
	correctResult := []FishingSpot{
		{1, &ap},
		{2, &rt},
		{3, &pi},
	}

	store := pgxsetup(t)
	defer pgxteardown(store, t)

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
