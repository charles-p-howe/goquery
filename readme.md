# goquery is a small library to simplify commodity database oprations.

## Warning: This library is in an experimental phase

---
Connecting to a RDBMS:

 - populate a rdbms configuration struct with the following info:
 ```go
   type RdbmsConfig struct {
	Dbuser      string 
	Dbpass      string
	Dbhost      string
	Dbport      string
	Dbname      string //db instance or name
	ExternalLib string //any external libraries required by the underlying db driver.  for example the instance client location for oracle connections
	DbDriver    string //db driver reference 
	DbStore     string //goquery store type.  Currently choices are 'pgx' or 'sqlx'
}
```
  - Driver Stores
    - pgx: uses the pgx driver and is postgres only
    - sqlx: uses sqlx and all sql compliant db drivers

<br/>
 Create the connection

 ```go
 store,err:=NewRdbmsDataStore(&config)
 ```

<br/>

---

## Querying
<br/>

- Minimal example into a slice
```go
dest:=[]mystruct{}
id:=10
err:=store.Select("select * from mytable where id>$id").
           Params(id).
	       Dest(&dest).
	       Fetch()
```

- Kitchen sink examples into a slice
```go
type MyFields struct{
	field1 string
	field2 float64
}

dest:=[]MyFields{}
id:=10
myTable:=data.TableDataset{
	Name:"mytable",
	Schema:"myschema",
	Statements:map[string]string{
		"select-all":`select * from mytable`,
		"select-none":`select * from mytable where false`,
	},
	Fields:MyFields{}, //depricated...do not use
}

err:=store.Select().
	Dataset(myTable).
	Dest(&dest).
	Fetch()

//or

err:=store.Select().
	Dataset(myTable).
	StatementKey("select-all")
	Dest(&dest).
	Fetch()

//or

err:=store.Select().
	Dataset(myTable).
	StatementKey("select-all").
	Suffix("where id<$1").
	Params(id).
	Dest(&dest).
	Fetch()

//or

err:=store.Select("select %s from %s").
	Suffix("where %s<$1").
	Appends("*","mytable","id"). //this is just string concatonation.  Never append user input. 
	Params(id).
	Dest(&dest).
	Fetch()

//finally you can query against transactions, and can optionally panic on err
err:=store.Select().
	Dataset(myTable).
	Tx(&tx) //transaction reference
	Dest(&dest).
	PanicOnErr(true).
	Fetch()

```

- As JSON
```go
id:=10
jsonb,err:=store.Select("select * from mytable where id>$id").
	Params(id).
	OmitNull(false). //will include null keys in output
	ForceArray(true). //will force output to a json array regardless of the number of records
	ToCamelCase(true). //will convert snake case to camel case
	DateFormat("02-Jan-2006"). //convert dates to formatted strings 
	FetchJSON()

```