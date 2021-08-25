### goquery is a small library to simplify commodity database oprations.

Creating an instance to a RDBMS:

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

 - create the connection

 ```go
 store,err:=NewRdbmsDataStore(&config)
 ```

 - query operations:
   - 
 ```go

 ```

