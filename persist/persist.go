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


type Record struct {
	FName string 
	FType string 
	FContent []byte 
	TimeStamp time.Time
}

// InitializeDatabase opens a database file and create the appropriate tables.
func ConnectDatabase(name string, dbpath string) (db *sql.DB, err error) {
	dbname := name + ".db"
	db, err = sql.Open("sqlite3", dbpath + dbname)
	if err != nil {
		// no such file
		log.Fatal(err)
	}
//	defer db.Close()
	sqlStmt := `
	create table if not exists runfiles (id integer not null primary key, filename text, filetype text, filecontent blob, timestamp DATETIME );
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
func InsertNewRecord(db *sql.DB, r Record) {
	// Check for existing file with the same file name.
	queryString := "select id, filename from runfiles where filename = " + "'" + r.FName + "'"
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
		stmt, err := db.Prepare("insert into runfiles(filename, filetype, filecontent, timestamp) values(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(r.FName, r.FType, r.FContent, r.TimeStamp)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// GetRecsByTime retrieves on or more binary blobs stored in the database for 
// a given date range provided as YYYY-MM-DD.
func GetRecsByTime(db *sql.DB, startTime, endTime time.Time) (recs []Record) {

	queryString := "select * from runfiles where timestamp >= '" + startTime.Format("2006-01-02 15:04:05")  + "' " + "and timestamp <= '" + endTime.Format("2006-01-02 15:04:05") + "'"
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	var result []Record
	for rows.Next() {
		var id int
		var fName, fType string
		var content []byte
		var tStamp time.Time
		err = rows.Scan(&id, &fName, &fType, &content, &tStamp )
		if err != nil {
			log.Fatal(err)
		}
		rec := Record{FName: fName, FType: fType, FContent: content, TimeStamp: tStamp }
		result = append(result, rec)
	}
	return result
}