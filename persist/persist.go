// Package persist saves run information between sessions of fitplot.
package persist

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	//"fmt"
	"time"
)


type Record struct {
	FName string 
	FType string 
	FContent []byte 
	TimeStamp time.Time
}

// MigrateDatabase updates the database schema to the latest version.
func MigrateDatabase(db *sql.DB) (err error) {
	// Setup the goose configuration
	migrateConf := &goose.DBConf{
		MigrationsDir: "./db/migrations",
		Env:           "production",
		Driver: goose.DBDriver{
			Name:    "sqlite3",
			OpenStr:  ".",
			Import:  "github.com/mattn/go-sqlite3",
			Dialect: &goose.Sqlite3Dialect{},
		},
	}
	// Get the latest possible migration
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		log.Printf("%q: %s\n", err, "Could not get latest upgrade!")
		return err
	}
	// Migrate up to the latest version
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db)
	if err != nil {
		log.Printf("%q: %s\n", err, "Could not migrate database!")
		return err
	}
	return nil
}
	

// InitializeDatabase opens a database file and create the appropriate tables.
func ConnectDatabase(name string, dbpath string) (db *sql.DB, err error) {
	dbname := name + ".db"
	// need to set a busy timeout (e.g. retry) when uploading multiple files to avoid locked db messages
	db, err = sql.Open("sqlite3", "file:" + dbpath + dbname + "?_busy_timeout=20000")
	if err != nil {
		// no such file or locked.
		log.Printf("%q: %s\n", err, "Could not open database! Locked?")
	}
	return db, err
}

// InsertNewRecord inserts a new record into the runfiles table containing a filename
// and a binary blob.  It assumes the database has been initialized and the table built.
func InsertNewRecord(db *sql.DB, r Record) (err error) {
	// Check for existing file with the same file name.
	queryString := "select id, filename from runfiles where filename = " + "'" + r.FName + "'"
	rows, err := db.Query(queryString)
	if err != nil {
		log.Printf("%q: %s\n", err, "Could not query database:" + r.FName)
		return err
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		found = true
	}
	// Insert a new row.
	if found == false {
		stmt, err := db.Prepare("insert into runfiles(filename, filetype, filecontent, timestamp) values(?,?,?,?)")
		if err != nil {
			log.Printf("%q: %s\n", err, "Could not prepare to insert into database!")
			return err
		}
		_, err = stmt.Exec(r.FName, r.FType, r.FContent, r.TimeStamp)
		if err != nil {
			log.Printf("%q: %s\n", err, "Could not insert into database!")
			return err
		}
	}
	return nil
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
			log.Printf("%q: %s\n", err, "Could not scan!")
		}
		rec := Record{FName: fName, FType: fType, FContent: content, TimeStamp: tStamp }
		result = append(result, rec)
	}
	return result
}

// CreateTempFile takes an in-memory array of bytes and stores it in a temporary
// file location (which varies by operating system)
func CreateTempFile(bytes []byte) (tmpFile *os.File, err error) {
	tmpFile, err = ioutil.TempFile("", "tmp")
	if err != nil {
		log.Printf("%q: %s %s\n", err, "Could not open", tmpFile.Name())
		return tmpFile, err
	}
	defer tmpFile.Close()
	tmpFile.Write(bytes)
	if err != nil {
		log.Printf("%q: %s\n", err, "Could not write to open temp file.")
		return nil, err
	}
	return tmpFile, nil
}