//
// Fitplot provides a webserver used to process .fit and .tcx files.
//
package main

import (
	//"database/sql"
	"encoding/csv"
//	"fmt"
	"github.com/cprevallet/fitplot/persist"
	"github.com/cprevallet/fitplot/tcx"
        "github.com/cprevallet/fitplot/strutil"
	"github.com/jezard/fit"
	"io/ioutil"
	"net/http"
	"os"
        "strconv"
)

var minPace = 0.56
var paceToMetric = 16.666667      // sec/meter -> min/km

func makePace( speed float64) ( pace float64) {
	// Speed -> pace
	if speed > 1.8 { //m/s
		pace = 1.0/speed
	} else {
		pace = minPace // s/m = 15 min/mi
	}
        return
        }

// Slice up a structure.
func makeRecs(runRecs []fit.Record) ( dumprec[][]string ) {
        // Calculate values from raw values.
	for _, record := range runRecs {
                newrec := []string{strconv.FormatFloat(record.Distance, 'G', -1, 64),
                                   strconv.FormatFloat(record.Speed, 'G', -1, 64),
                                   strutil.DecimalTimetoMinSec(makePace(record.Speed) * paceToMetric),
                                   strconv.FormatFloat(record.Position_lat, 'G', -1, 64),
                                   strconv.FormatFloat(record.Position_long, 'G', -1, 64),
                                   strconv.FormatFloat(record.Altitude, 'G', -1, 64),
                                   strconv.Itoa(int(record.Cadence))}
                dumprec = append(dumprec, newrec)
        }
		return dumprec
	}

// Parse the input bytes into structures more conducive for additional
// processing by routines in graphdata.go.
func parseInputBytes(fBytes []byte) (fType string, fitStruct fit.FitFile, tcxdb *tcx.TCXDB, runRecs []fit.Record, runLaps []fit.Lap, err error) {
	// Make a copy in a temporary folder for use with fit and tcx
	// libraries.
	tmpFile, err := persist.CreateTempFile(fBytes)
	if err != nil {
		return "", fit.FitFile{}, nil, nil, nil, err
	}
	err = nil
	tmpFname := tmpFile.Name()
	// Determine what type of file we're looking at.
	rslt := http.DetectContentType(fBytes)
	switch {
	case rslt == "application/octet-stream":
		// Filetype is FIT, or at least it could be?
		fType = "FIT"
		fitStruct = fit.Parse(tmpFname, false)
		tcxdb = nil
		runRecs = fitStruct.Records
		runLaps = fitStruct.Laps
	case rslt == "text/xml; charset=utf-8":
		// Filetype is TCX or at least it could be?
		fType = "TCX"
		fitStruct = fit.FitFile{}
		tcxdb, err = tcx.ReadTCXFile(tmpFname)
		// We cleverly convert the values of interest into a structures we already
		// can handle.
		if err == nil {
			runRecs = tcx.CvtToFitRecs(tcxdb)
			runLaps = tcx.CvtToFitLaps(tcxdb)
		}
	}
	persist.DeleteTempFile(tmpFile)
	return fType, fitStruct, tcxdb, runRecs, runLaps, err
}

// Parse the uploaded file, parse it and return run information suitable
// to construct the user interface.
func dumpCSV(fBytes []byte) {
	// Parse the input bytes into structures more conducive for additional
	// processing by routines in graphdata.go.
	_, _, _, runRecs, _, err := parseInputBytes(fBytes)
	if err != nil {
                panic(err)
	}
        outrecs := makeRecs(runRecs)
	w := csv.NewWriter(os.Stdout)
	w.WriteAll(outrecs)
	if err := w.Error(); err != nil {
		panic(err)
	}
}

func main() {
	//openLog()
        file, err := os.Open("INPUTDATA.FIT") // For read access.
	if err != nil {
		panic("Can't open log file. Aborting.")
	}
	// fBytes is an in-memory array of bytes read from the file.
	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Can't open input file. Aborting.")
	}
        dumpCSV(fBytes)
        file.Close()
}
