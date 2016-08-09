// Package persist saves run information between sessions of fitplot.
package persist

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//"os"
	//"fmt"
)

// InitializeDatabase opens a database file and create the appropriate tables.
// WARNING: Existing database of the same name will be DELETED.
func ConnectDatabase(name string, dbpath string) (db *sql.DB, err error) {
	_ = "breakpoint"
	dbname := name + ".db"
	//finfo, err := os.Stat(dbpath + dbname)
	db, err = sql.Open("sqlite3", dbpath + dbname)
	if err != nil {
		// no such file
		log.Fatal(err)
	}
//	defer db.Close()
	sqlStmt := `
	create table if not exists runfiles (id integer not null primary key, filename text, filetype text, content blob, timestamp text );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	return db, err
}

// InsertNewRecord inserts a new record into the runfiles table containing a filename
// and a binary blob.
func InsertNewRecord(db *sql.DB, fName string, fType string, content []byte, timestamp string) {
	// insert
	stmt, err := db.Prepare("insert into runfiles(filename, filetype, content, timestamp) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	// TODO need to figure out how to retrieve last id entered.
	_, err = stmt.Exec(fName, fType, content, timestamp)
	if err != nil {
		log.Fatal(err)
	}
}
