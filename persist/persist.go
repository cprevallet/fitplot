// Package persist saves run information between sessions of fitplot.
package persist

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//"os"
	//"fmt"
	"time"
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
// and a binary blob.  It assumes the database has been initialized and the table built.
func InsertNewRecord(db *sql.DB, fName string, fType string, content []byte, timestamp time.Time) {
	// Check for existing file with the same file name.
	queryString := "select id, filename from runfiles where filename = " + "'" + fName + "'"
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	found := false
	for rows.Next() {
		found = true
	}
	// Insert a new row.
	if found == false {
		stmt, err := db.Prepare("insert into runfiles(filename, filetype, content, timestamp) values(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(fName, fType, content, timestamp)
		if err != nil {
			log.Fatal(err)
		}
	}
}
