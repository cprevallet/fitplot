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

// GetFileByTimeStamp retrieves on or more binary blobs stored in the database for 
// a given day provided by a timestamp.
func GetFileByTimeStamp(db *sql.DB, timestamp time.Time) (recs []Record) {

	thisDay := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
	thisDate := thisDay.Format("2006-01-02")
	addOneDay := timestamp.Add(time.Hour * 24)
	nextDay := time.Date(addOneDay.Year(), addOneDay.Month(), addOneDay.Day(), 0, 0, 0, 0, time.UTC)
	nextDate := nextDay.Format("2006-01-02")
	
	// Between is inclusive.
	queryString := "select * from runfiles where timestamp between '" + thisDate + "' " + "and '" + nextDate + "'"
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