//
// Fitplot provides a webserver used to process .fit and .tcx files.
//
package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cprevallet/fitplot/desktop"
	"github.com/cprevallet/fitplot/persist"
	"github.com/cprevallet/fitplot/strutil"
	"github.com/cprevallet/fitplot/tcx"
	"github.com/jezard/fit"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// Buildstamp represents a build timestamp for use in technical support.  It is
// created using the linker option -X we can set a value for a symbol that can
// be accessed from within the binary.
// go build -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" fitplot.go
var Buildstamp = "No build timestamp provided"

// Githash represents a hash from the version control system for use in
// technical support.  It is created using the linker option -X we can
// set a value for a symbol that can be accessed from within the binary.
// go build -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" fitplot.go
var Githash = "No git hash provided"
var tmpFname = ""
var timeStamp time.Time

// Compile templates on start for better performance.
var templates = template.Must(template.ParseFiles("tmpl/fitplot.html"))

// Display the named template.
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

// Handle requests to "/".
func pageloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Load page.
		//fmt.Println("pageloadHandler Received Request")
		display(w, "fitplot", nil)
	}
	if r.Method == "POST" {
		// File load request
		//fmt.Println("pageloadHandler POST Received Request")
		uploadHandler(w, r)
		//display(w, "fitplot", nil)
	}
}

//Upload a copy the fit file to a temporary local directory.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("uploadHandler Received Request")
	// Do some low-level stuff to retrieve the file name and byte array.
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fBytes is an in-memory array of bytes read from the file.
	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fName := handler.Filename
	fType, fitStruct, tcxdb, _, _, err := parseInputBytes(fBytes)
	switch {
	case fType == "FIT":
		timeStamp = time.Unix(fitStruct.Records[0].Timestamp, 0)
	case fType == "TCX":
		timeStamp = time.Unix(tcxdb.Acts.Act[0].Laps[0].Trk.Pt[0].Time.Unix(), 0)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Persist the in-memory array of bytes to the database.
	dbRecord := persist.Record{FName: fName, FType: fType, FContent: fBytes, TimeStamp: timeStamp}
	db, err := persist.ConnectDatabase("fitplot", "./")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = persist.InsertNewRecord(db, dbRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		db.Close()
		return
	}
	db.Close()
}

// getDistance retrieves a run's total distance in the appropriate unit system.
func getOtherVals(fBytes []byte) (totalDistance float64, movingTime float64, totalPace string) {
	dummyTimeStamp := []int64{0}
	lapCal := []float64{0.0}
	_, _, _, _, runLaps, _ := parseInputBytes(fBytes)
	_, _, _, _, totalDistance, movingTime = processFitLap(runLaps, toEnglish)
	_,  totalPace, _, _, _, _,_ = createStats(toEnglish, totalDistance, movingTime, dummyTimeStamp, lapCal)
	return totalDistance, movingTime, totalPace
}


// Return information about entries in the database between two dates.
func dbGetRecs(w http.ResponseWriter, r *http.Request) (recs []persist.Record, err error) {
	type DBDateStrings struct {
		DBStart string
		DBEnd   string
	}
	decoder := json.NewDecoder(r.Body)
	var dbQuery DBDateStrings //string
	err = decoder.Decode(&dbQuery)
	if err != nil {
		return nil, err
	}
	// Connect to database and retrieve.  Be sure we bracket the entirety of
	// the selected day.
	db, err := persist.ConnectDatabase("fitplot", "./")
	if err != nil {
		return nil, err
	}
	startTime, _ := time.Parse("2006-01-02 15:04:05", dbQuery.DBStart+" 00:00:00")
	endTime, _ := time.Parse("2006-01-02 15:04:05", dbQuery.DBEnd+" 23:59:59")
	recs = persist.GetRecsByTime(db, startTime, endTime)
	db.Close()
	return recs, nil
}

// Return information about entries in the database.
func dbHandler(w http.ResponseWriter, r *http.Request) {
	type RunInfoStruct struct {
		FName		string
		FType		string
		TimeStamp	string
		Date		string
		Time		string
		TimeZone	string
		Weekday		string
		MovingTime	string
		Pace        string
		Distance    string
	}
	var DBFileList []RunInfoStruct
	// Structure element names MUST be uppercase or decoder can't access them.
	type RtnStruct struct {
		DBFileList []RunInfoStruct
		Totals map[string]float64
		Units map[string]string
	}
	returnData := RtnStruct{
		DBFileList: nil,
		Totals:     nil,
		Units:      nil,
	}
	totals := map[string]float64 {"Distance": 0.0}
    recs, err := dbGetRecs(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Do a bit of fancy footwork to speed up the reads with goroutines.
	ch := make(chan RunInfoStruct, 100)  //make this buffered only 100 at a time
	for _, rec := range recs {
		go func(rec persist.Record) {
			// process
			var rs RunInfoStruct
			rs.FName = rec.FName
			rs.FType = rec.FType
			rs.TimeStamp = rec.TimeStamp.Format(time.RFC1123)
			rs.Date = rec.TimeStamp.Format(time.RFC3339)[0:10]
			rs.Time = rec.TimeStamp.Format(time.RFC3339)[11:19]
			rs.TimeZone = rec.TimeStamp.Format(time.RFC3339)[19:25]
			rs.Weekday = rec.TimeStamp.Format(time.RFC1123)[0:3]
			totalDistance, movingTime,totalPace := getOtherVals(rec.FContent)
			rs.Distance = strconv.FormatFloat(totalDistance, 'f', 2, 64)
			rs.MovingTime = strutil.DecimalTimetoHourMinSec(movingTime)
			rs.Pace = totalPace
			ch <- rs
		}(rec)
		rs1 := <-ch
		f, _ := strconv.ParseFloat(rs1.Distance, 64)
		totals["Distance"] += f
		DBFileList = append(DBFileList, rs1)
	}
	var units map[string]string
	units = make(map[string]string)
	if toEnglish {
		units["Distance"] = "miles"
	} else {
		units["Distance"] = "kilometers"
	}
	returnData.DBFileList = DBFileList
	returnData.Totals = totals
	returnData.Units = units
	
	//Convert to json.
	js, err := json.Marshal(returnData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Send
	w.Write(js)
}

// Set an individually record in the the database as the one to process via
// the timeStamp global variable.
func dbSelectHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type DateStrings struct {
		DateStr string
		TimeStr string
		TimeZoneStr string
	}
	var ds DateStrings
	var ts string
	err := decoder.Decode(&ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ts = ds.DateStr + "T"+ ds.TimeStr + ds.TimeZoneStr
	timeStamp, err = time.Parse(time.RFC3339, ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Export file(s) to disk.
func dbExportHandler(w http.ResponseWriter, r *http.Request) {
    recs,err := dbGetRecs(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	slash := string(filepath.Separator)
	path := "." + slash + "export" + slash
	_ = "breakpoint"
	for _, rec := range recs {
		ioutil.WriteFile(path + rec.FName, rec.FContent, 0644)
	}
}

// Return information about the runtime environment.
func envHandler(w http.ResponseWriter, r *http.Request) {

	type Environ struct {
		Buildstamp      string
		Githash         string
		CPUArchitecture string
		OperatingSystem string
	}

	e := Environ{
		Buildstamp:      Buildstamp,
		Githash:         Githash,
		CPUArchitecture: runtime.GOARCH,
		OperatingSystem: runtime.GOOS,
	}

	//Convert to json.
	js, err := json.Marshal(e)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println("plotHandler Received Request")
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Send
	w.Write(js)
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
	tmpFname = tmpFile.Name()
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
func plotHandler(w http.ResponseWriter, r *http.Request) {
	type Plotvals struct {
		Titletext         string
		XName             string
		Y0Name            string
		Y1Name            string
		Y2Name            string
		DispDistance      []float64
		DispPace          []float64
		DispAltitude      []float64
		DispCadence       []float64
		TimeStamps        []int64
		Latlongs          []map[string]float64
		LapDist           []float64
		LapTime           []string
		LapCal            []float64
		LapPace           []string
		C0Str             string
		C1Str             string
		C2Str             string
		C3Str             string
		C4Str             string
		TotalDistance     float64
		MovingTime        float64
		DispTotalDistance string
		DispMovingTime    string
		TotalPace         string
		ElapsedTime       string
		TotalCal          string
		AvgPower          string
		StartDateStamp    string
		EndDateStamp      string
		Device            string
		PredictedTimes    map[string]string
		TrainingPaces     map[string]string
		VDOT              float64
		DeviceName        string
		DeviceUnitID      string
		DeviceProdID      string
		RunScore          float64
		VO2max            float64
		Buildstamp        string
		Githash           string
	}
	// Note extra space on following assignments.
	var xStr = "Distance "
	var y0Str = "Pace "
	var y1Str = "Elevation "
	var y2Str = "Cadence "
	var c0Str, c1Str, c2Str, c3Str, c4Str string

	// User hasn't selected a file yet?  Avoid a panic.
	if timeStamp.IsZero() {
		http.Error(w, "No file selected.", http.StatusConflict)
		return
	}
	// Structure element names MUST be uppercase or decoder can't access them.
	type UIData struct {
		UseEnglish bool
		Racedist   float64
		Racehours  int64
		Racemins   int64
		Racesecs   int64
		UseSegment bool
		Splitdist  float64
		Splithours int64
		Splitmins  int64
		Splitsecs  int64
	}
	decoder := json.NewDecoder(r.Body)
	var UI UIData
	err := decoder.Decode(&UI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	toEnglish = UI.UseEnglish

	// Retrieve the file from the database by timeStamp global
	// variable.  Make the search criteria just outside the
	// expected run start time.
	db, err := persist.ConnectDatabase("fitplot", "./")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slightlyOlder := timeStamp.Add(-1 * time.Second)
	slightlyNewer := timeStamp.Add(1 * time.Second)
	recs := persist.GetRecsByTime(db, slightlyOlder, slightlyNewer)
	db.Close()
	fBytes := recs[0].FContent

	// Parse the input bytes into structures more conducive for additional
	// processing by routines in graphdata.go.
	fType, fitStruct, tcxdb, runRecs, runLaps, err := parseInputBytes(fBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build the variable strings based on unit system.
	if toEnglish {
		xStr = xStr + "(mi)"
		y0Str = y0Str + "(min/mi)"
		y1Str = y1Str + "(ft)"
		y2Str = y2Str + "(steps/min)"
		c0Str = "Lap"
		c1Str = "Distance" + "(mi)"
		c2Str = "Pace" + "(min/mi)"
		c3Str = "Time" + "(min)"
		c4Str = "Calories" + "(kcal)"
	} else {
		xStr = xStr + "(km)"
		y0Str = y0Str + "(min/km)"
		y1Str = y1Str + "(m)"
		y2Str = y2Str + "(steps/min)"
		c0Str = "Lap"
		c1Str = "Distance" + "(km)"
		c2Str = "Pace" + "(min/km)"
		c3Str = "Time" + "(min)"
		c4Str = "Calories" + "(kcal)"
	}

	// Create an object to contain various plot values.
	p := Plotvals{Titletext: "",
		XName:             xStr,
		Y0Name:            y0Str,
		Y1Name:            y1Str,
		Y2Name:            y2Str,
		DispDistance:      nil,
		DispPace:          nil,
		DispAltitude:      nil,
		DispCadence:       nil,
		TimeStamps:        nil,
		Latlongs:          nil,
		LapDist:           nil,
		LapTime:           nil,
		LapCal:            nil,
		LapPace:           nil,
		C0Str:             c0Str,
		C1Str:             c1Str,
		C2Str:             c2Str,
		C3Str:             c3Str,
		C4Str:             c4Str,
		TotalDistance:     0.0,
		MovingTime:        0.0,
		DispTotalDistance: "",
		DispMovingTime:    "",
		TotalPace:         "",
		TotalCal:          "",
		AvgPower:          "",
		StartDateStamp:    "",
		EndDateStamp:      "",
		Device:            "",
		PredictedTimes:    nil,
		TrainingPaces:     nil,
		VDOT:              0.0,
		DeviceName:        "",
		DeviceUnitID:      "",
		DeviceProdID:      "",
		RunScore:          0.0,
		VO2max:            0.0,
	}

	// Retrieve overview information.
	if fType == "FIT" {
		p.DeviceName = fitStruct.DeviceInfo[0].Manufacturer
		p.DeviceProdID = fitStruct.DeviceInfo[0].Product
		p.DeviceUnitID = fmt.Sprint(fitStruct.DeviceInfo[0].Serial_number)
	}
	if fType == "TCX" {
		p.DeviceName, p.DeviceUnitID, p.DeviceProdID = tcx.DeviceInfo(tcxdb)
	}

	// Here's where the heavy lifting of pulling tracks and performance information
	// from (portions of) the fit file into something we can view is done.
	p.Latlongs, p.TimeStamps, p.DispDistance, p.DispPace, p.DispAltitude, p.DispCadence =
		processFitRecord(runRecs, toEnglish)

	p.LapDist, p.LapTime, p.LapCal, p.LapPace, p.TotalDistance, p.MovingTime = processFitLap(runLaps, toEnglish)

	// Get the start time.
	p.Titletext += time.Unix(p.TimeStamps[0], 0).Format(time.UnixDate)

	// Calculate the summary string information.
	p.DispTotalDistance,  p.TotalPace, p.DispMovingTime, p.TotalCal, p.AvgPower, p.StartDateStamp,
		p.EndDateStamp = createStats(toEnglish, p.TotalDistance, p.MovingTime, p.TimeStamps, p.LapCal)

	// Calculate the analysis page.
	p.PredictedTimes, p.VDOT, p.VO2max, p.RunScore, p.TrainingPaces = createAnalysis(toEnglish,
		UI.UseSegment, p.DispDistance, p.TimeStamps, UI.Splitdist, UI.Splithours, UI.Splitmins,
		UI.Splitsecs, UI.Racedist, UI.Racehours, UI.Racemins, UI.Racesecs)

	//Convert to json.
	js, err := json.Marshal(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println("plotHandler Received Request")
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Send
	w.Write(js)
}

// openLog opens a log file.
func openLog(){
	f, err := os.OpenFile("fitplot.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic("Can't open log file. Aborting.")
	}
	log.SetOutput(f)
}

// Allow the user to shutdown the server from a target in the user interface client.
func stopHandler(w http.ResponseWriter, r *http.Request) {
	// Nothing graceful about this exit.  Just bail out.
	os.Exit(0)
	return
}

// Migrate the persistent store (database) to the current version.
func migrate() {
	// Play it safe and create a single backup file in the temporary 
	// directory in the event of problems.
	// I don't know what the performance impact might be, this safeguard
	// can be removed if performance becomes an issue (e.g. for large databases).
	file, err := os.Open("fitplot.db")
	if err == nil {
		fBytes,_ := ioutil.ReadAll(file)
		persist.CreateTempFile(fBytes)
		file.Close()
	}
	db, err := persist.ConnectDatabase("fitplot", "./")
	if err != nil {
		panic("Can't locate database. Aborting.")
	}
	persist.MigrateDatabase(db)
	db.Close()
	return
}

func main() {
	openLog()
	desktop.Open("http://localhost:8080")
	// Serve static files if the prefix is "static".
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	migrate()
	// Handle normal requests.
	http.HandleFunc("/", pageloadHandler)
	http.HandleFunc("/getplot", plotHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/getruns", dbHandler)
	http.HandleFunc("/selectrun", dbSelectHandler)
	http.HandleFunc("/exportrun", dbExportHandler)
	//Listen on port 8080
	//fmt.Println("Server starting on port 8080.")
	http.ListenAndServe(":8080", nil)
}
