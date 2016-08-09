// Package persist saves run information between sessions of fitplot.
package persist

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

// InitializeDatabase opens a database file and create the appropriate tables.
// WARNING: Existing database of the same name will be DELETED.
func InitializeDatabase(name string, dbpath string) (db *sql.DB, err error) {
	dbname := name + ".db"
	os.Remove(dbpath + dbname)
	
	db, err = sql.Open("sqlite3", dbpath + dbname)
	if err != nil {
		log.Fatal(err)
	}
//	defer db.Close()
	
	sqlStmt := `
	create table runfiles (id integer not null primary key, filename text, file blob);
	delete from runfiles;
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
func InsertNewRecord(db *sql.DB, filename string, file []byte) {
	// insert
	stmt, err := db.Prepare("insert into runfiles(id, filename, file) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	// TODO need to figure out how to retrieve last id entered.
	_, err = stmt.Exec(1, filename, file)
	if err != nil {
		log.Fatal(err)
	}
}
