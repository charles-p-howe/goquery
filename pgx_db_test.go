package goquery

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type FishingSpot struct {
	ID       int32   `db:"id"`
	Location *string `db:"location"`
}

func pgxsetup(t *testing.T) DataStore {
	ctx := context.Background()
	store := getPgxStore(t)
	err := store.Transaction(func(tx Tx) {
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
		_, err = tx.PgxTx().Exec(ctx, "create table json_test (id int,json_attr jsonb)")
		if err != nil {
			panic(err)
		}

		_, err = tx.PgxTx().Exec(ctx, "create table arrays_test (id int,intlist integer[],boollist boolean[])")
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		t.Errorf("Setup error:%s\n", err)
	}
	return store
}

func pgxteardown(store DataStore, t *testing.T) {
	//ctx := context.Background()
	err := store.Transaction(func(tx Tx) {
		store.MustExec(&tx, "drop table fishing_spots")
		store.MustExec(&tx, "drop table json_test")
		store.MustExec(&tx, "drop table arrays_test")
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
	store := RdbmsDataStore{&db}
	return &store
}

func TestPgxConnection(t *testing.T) {
	getPgxStore(t)
}

func TestPgxJsonDepricated(t *testing.T) {
	correctResult := `[{"id":1,"location":"Alpine Frove"},{"id":2,"location":"Rivertown"},{"id":3,"location":"Pine Island"},{"id":4,"location":null}]`
	store := pgxsetup(t)
	defer pgxteardown(store, t)

	json, err := store.
		Select("select * from fishing_spots").
		OmitNull(false).
		IsJsonArray(true).
		FetchJSON()
	if err != nil {
		t.Errorf("Failed JSON Test: %s\n", err)
	}
	jsonstring := string(json)
	if jsonstring != correctResult {
		t.Errorf("Failed JSON Test: Got %s want %s", jsonstring, correctResult)
	}
}

func TestPgxJson(t *testing.T) {
	correctResult := `[{"id":1,"location":"Alpine Frove"},{"id":2,"location":"Rivertown"},{"id":3,"location":"Pine Island"},{"id":4,"location":null}]`
	store := pgxsetup(t)
	defer pgxteardown(store, t)

	builder := strings.Builder{}
	err := store.
		Select("select * from fishing_spots").
		OmitNull(false).
		IsJsonArray(true).
		OutputJson(&builder).
		Fetch()
	if err != nil {
		t.Errorf("Failed JSON Test: %s\n", err)
	}
	jsonstring := builder.String()
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
		{4, nil},
	}

	store := pgxsetup(t)
	defer pgxteardown(store, t)

	///////////autogenerate select///////////////////
	fsTbl := TableDataSet{
		Name: "fishing_spots",
	}

	dest := &[]FishingSpot{}
	err := store.Select().
		DataSet(&fsTbl).
		Dest(dest).
		PanicOnErr(true).
		Fetch()
	if err != nil {
		t.Errorf("Failed Slice Test:%s\n", err)
	}

	if reflect.DeepEqual(dest, correctResult) {
		t.Errorf("Failed Slice Test: Got %v want %v", dest, correctResult)
	}

	/////////////////add select statement///////////
	stmts := Statements{
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

func TestPgxRow(t *testing.T) {
	store := pgxsetup(t)
	defer pgxteardown(store, t)
	stmt := "select * from fishing_spots"
	dest := FishingSpot{}
	rows, err := store.Select(stmt).FetchRows()
	if err != nil {
		t.Errorf("Failed Rows Test:%s\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.ScanStruct(&dest)
		fmt.Println(dest)
	}
}

func TestPgxRowFunction(t *testing.T) {
	store := pgxsetup(t)
	defer pgxteardown(store, t)
	stmt := "select * from fishing_spots"
	dest := FishingSpot{}
	builder := strings.Builder{}
	err := store.Select(stmt).
		ForEachRow(func(row Rows) error {
			if dest.Location != nil {
				builder.WriteString(fmt.Sprintf("%d:%s\n", dest.ID, *dest.Location))
			}
			return nil
		}).
		Dest(&dest).
		Fetch()
	if err != nil {
		t.Errorf("Failed Rows Test:%s\n", err)
	}
	fmt.Println(builder.String())
}

func TestPgxRowFunction2(t *testing.T) {
	store := pgxsetup(t)
	defer pgxteardown(store, t)
	stmt := "select * from fishing_spots"
	dest := FishingSpot{}
	builder := strings.Builder{}
	err := store.Select(stmt).
		ForEachRow(func(row Rows) error {
			err := row.ScanStruct(&dest)
			if err != nil {
				return err
			}
			if dest.Location != nil {
				builder.WriteString(fmt.Sprintf("%d:%s\n", dest.ID, *dest.Location))
			}
			return nil
		}).
		Dest(&dest).
		Fetch()
	if err != nil {
		t.Errorf("Failed Rows Test:%s\n", err)
	}
	fmt.Println(builder.String())
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
		Name:        "fishing_spots",
		TableFields: FishingSpot{},
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
		Name:        "fishing_spots",
		TableFields: FishingSpot{},
	}

	store := pgxsetup(t)
	defer pgxteardown(store, t)
	err := store.Insert(&fsTbl).Records(&fs).Batch(true).BatchSize(len(fs)).Execute()
	if err != nil {
		t.Error(err)
	}

}

type ArrayTest struct {
	ID       int    `db:"id"`
	IntList  []int  `db:"intlist"`
	Boollist []bool `db:"boollist"`
}

var at TableDataSet = TableDataSet{
	Name:        "arrays_test",
	TableFields: ArrayTest{},
}

func TestPgxArrayInsert(t *testing.T) {
	store := pgxsetup(t)
	defer pgxteardown(store, t)
	il := []int{3, 5, 7, 9}
	bl := []bool{false, true, false, true}
	store.MustExec(NoTx, `insert into arrays_test (id,intlist,boollist) values ($1,$2,$3)`, 1, il, bl)

	atData := ArrayTest{2, []int{100, 101, 102}, []bool{true, false, false, true}}
	err := store.Insert(&at).Records(atData).Execute()
	if err != nil {
		t.Error(err)
	}

	dest := []ArrayTest{}
	err = store.Select("select * from arrays_test").Dest(&dest).Fetch()
	if err != nil {
		t.Error(err)
	}
	t.Log(dest)
}

type JsonTest struct {
	ID         int      `db:"id"`
	Attributes JsonAttr `db:"json_attr"`
}

type JsonAttr struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var fs TableDataSet = TableDataSet{
	Name:        "json_test",
	TableFields: JsonTest{},
}

func TestJson(t *testing.T) {
	store := pgxsetup(t)
	defer pgxteardown(store, t)
	store.MustExec(NoTx, `insert into json_test values (1,'{"name":"jack","age":8}')`)
	store.MustExec(NoTx, `insert into json_test values ($1,$2)`, 2, JsonAttr{"Luna", 4})
	recs := []JsonTest{
		{3, JsonAttr{"John", 20}},
		{4, JsonAttr{"Karen", 30}},
	}
	err := store.Insert(&fs).Records(recs).Execute()
	if err != nil {
		t.Error(err)
	}
	//store.MustExec(NoTx, `insert into json_test (id) values (1,'{"n`)
	dest := []JsonTest{}
	err = store.Select("select * from json_test").Dest(&dest).Fetch()
	if err != nil {
		t.Error(err)
	}
	t.Log(dest)
}
