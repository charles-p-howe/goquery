package goquery

import (
	"context"
	"reflect"
	"strconv"
	"testing"
)

type FishingSpot struct {
	ID       int32   `db:"id"`
	Location *string `db:"location"`
}

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
	//dialect, _ := getDialect("pgx")
	db, err := NewPgxConnection(config)
	if err != nil {
		t.Errorf("Failed to connect to store:%s\n", err)
	}

	store := RdbmsDataStore{&db}
	return &store
}

func TestPgxConnection(t *testing.T) {
	getPgxStore(t)
}

func TestPgxJson(t *testing.T) {
	correctResult := `[{"id":1,"location":"Alpine Frove"},{"id":2,"location":"Rivertown"},{"id":3,"location":"Pine Island"},{"id":4,"location":null}]`
	store := pgxsetup(t)
	defer pgxteardown(store, t)

	json, err := store.
		Select("select * from fishing_spots").
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
	fsTbl := TableDataSet{
		Name:   "fishing_spots",
		Fields: FishingSpot{},
	}

	dest := &[]FishingSpot{}
	err := store.Select().DataSet(&fsTbl).Dest(dest).PanicOnErr(true).Fetch()
	if err != nil {
		t.Errorf("Failed Slice Test:%s\n", err)
	}

	if reflect.DeepEqual(dest, correctResult) {
		t.Errorf("Failed Slice Test: Got %v want %v", dest, correctResult)
	}

	/////////////////add select statement///////////
	stmts := map[string]string{
		"named-select": `select * from fishing_spots`,
	}

	fsTbl.Statements = stmts
	dest = &[]FishingSpot{}
	err = store.Select().
		DataSet(&fsTbl).
		StatementKey("named-select").
		Dest(dest).
		Fetch()

	if err != nil {
		t.Errorf("Failed Slice Test:%s\n", err)
	}

	if !reflect.DeepEqual(dest, correctResult) {
		t.Errorf("Failed Slice Test: Got %v want %v", dest, correctResult)
	}

}

func TestPgxInsert(t *testing.T) {
	l10 := "New Spot 10"
	l11 := "New Spot 11"
	fs := []FishingSpot{
		{10, &l10},
		{11, &l11},
	}

	//fs := FishingSpot{10, &l10}

	fsTbl := TableDataSet{
		Name:   "fishing_spots",
		Fields: FishingSpot{},
	}

	store := pgxsetup(t)
	defer pgxteardown(store, t)
	err := store.Insert(&fsTbl).Records(fs).Execute()
	if err != nil {
		t.Error(err)
	}

}

func TestPgxInsertBatch(t *testing.T) {

	/*
		l10 := "New Spot 10"
		l11 := "New Spot 11"
		fs := []FishingSpot{
			{10, &l10},
			{11, &l11},
		}
	*/

	count := 40000
	fs := make([]FishingSpot, count)
	for i := 0; i < count; i++ {
		val := strconv.Itoa(i)
		fs[i] = FishingSpot{int32(i + 10), &val}
	}

	fsTbl := TableDataSet{
		Name:   "fishing_spots",
		Fields: FishingSpot{},
	}

	store := pgxsetup(t)
	defer pgxteardown(store, t)
	err := store.Insert(&fsTbl).Records(&fs).Batch(true).BatchSize(len(fs)).Execute()
	if err != nil {
		t.Error(err)
	}

}
