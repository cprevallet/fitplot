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
	"github.com/cprevallet/fitplot/tcx"
	"github.com/jezard/fit"
	"html/template"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"runtime"
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
			timeStamp = time.Unix(tcxdb.Acts.Act[0].Laps[0].Trk.Pt[0].Time.Unix(),0)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Persist the in-memory array of bytes to the database.
	dbRecord := persist.Record{FName: fName, FType: fType, FContent: fBytes, TimeStamp: timeStamp}
	db, _ := persist.ConnectDatabase("fitplot", "./")
	persist.InsertNewRecord(db, dbRecord)
	db.Close()
}

// Return information about entries in the database .
func dbHandler(w http.ResponseWriter, r *http.Request) {
	// Structure element names MUST be uppercase or decoder can't access them.
	type DBDateStrings struct {
		DBStart      string
		DBEnd        string
	}
	
	var DBFileList []map[string]string
	
	decoder := json.NewDecoder(r.Body)
	var dbQuery DBDateStrings //string
	err := decoder.Decode(&dbQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Connect to database and retrieve.  Be sure we bracket the entirety of
	// the selected day.
	db, _ := persist.ConnectDatabase("fitplot", "./")
	startTime, _ := time.Parse("2006-01-02 15:04:05", dbQuery.DBStart + " 00:00:00")
	endTime, _ := time.Parse("2006-01-02 15:04:05", dbQuery.DBEnd + " 23:59:59")
	recs := persist.GetRecsByTime(db, startTime, endTime )
	db.Close()
	for _, rec := range recs {
		var filerec map[string]string
		filerec = make(map[string]string)
		filerec["File name"] = rec.FName
		filerec["File type"] = rec.FType
		filerec["Timestamp"] = rec.TimeStamp.Format(time.RFC1123)
		DBFileList = append(DBFileList, filerec)
	}
	//Convert to json.
	js, err := json.Marshal(DBFileList)
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
	var ts string
	err := decoder.Decode(&ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// ts returned as RFC1123 string 
	timeStamp, err = time.Parse(time.RFC1123Z, ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func parseInputBytes (fBytes []byte)(fType string, fitStruct fit.FitFile, tcxdb *tcx.TCXDB, runRecs []fit.Record, runLaps []fit.Lap, err error) {
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
		tcxdb, err := tcx.ReadTCXFile(tmpFname)
		// We cleverly convert the values of interest into a structures we already
		// can handle.
		if err != nil {
			runRecs = tcx.CvtToFitRecs(tcxdb)
			runLaps = tcx.CvtToFitLaps(tcxdb)
		}
	}
	return fType, fitStruct, tcxdb, runRecs, runLaps, err
}

// Parse the uploaded file, parse it and return run information suitable
// to construct the user interface.
func plotHandler(w http.ResponseWriter, r *http.Request) {
	type Plotvals struct {
		Titletext      string
		XName          string
		Y0Name         string
		Y1Name         string
		Y2Name         string
		DispDistance   []float64
		DispPace       []float64
		DispAltitude   []float64
		DispCadence    []float64
		TimeStamps     []int64
		Latlongs       []map[string]float64
		LapDist        []float64
		LapTime        []string
		LapCal         []float64
		LapPace        []string
		C0Str          string
		C1Str          string
		C2Str          string
		C3Str          string
		C4Str          string
		TotalDistance  float64
		DispTotalDistance  string
		TotalPace      string
		ElapsedTime    string
		TotalCal       string
		AvgPower       string
		StartDateStamp string
		EndDateStamp   string
		Device         string
		PredictedTimes map[string]string
		TrainingPaces  map[string]string
		VDOT           float64
		DeviceName     string
		DeviceUnitID   string
		DeviceProdID   string
		RunScore       float64
		VO2max         float64
		Buildstamp     string
		Githash        string
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
		UseEnglish   bool
		Racedist 	 float64
		Racehours	 int64
		Racemins	 int64
		Racesecs	 int64
		UseSegment	 bool
		Splitdist	 float64
		Splithours	 int64
		Splitmins	 int64
		Splitsecs	 int64
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
	db, _ := persist.ConnectDatabase("fitplot", "./")
	slightlyOlder := timeStamp.Add(-1 * time.Second)
	slightlyNewer := timeStamp.Add(1 * time.Second)
	recs := persist.GetRecsByTime(db, slightlyOlder, slightlyNewer )
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
		XName:          xStr,
		Y0Name:         y0Str,
		Y1Name:         y1Str,
		Y2Name:         y2Str,
		DispDistance:   nil,
		DispPace:       nil,
		DispAltitude:   nil,
		DispCadence:    nil,
		TimeStamps:     nil,
		Latlongs:       nil,
		LapDist:        nil,
		LapTime:        nil,
		LapCal:         nil,
		LapPace:        nil,
		C0Str:          c0Str,
		C1Str:          c1Str,
		C2Str:          c2Str,
		C3Str:          c3Str,
		C4Str:          c4Str,
		TotalDistance:  0.0,
		DispTotalDistance:  "",
		TotalPace:      "",
		ElapsedTime:    "",
		TotalCal:       "",
		AvgPower:       "",
		StartDateStamp: "",
		EndDateStamp:   "",
		Device:         "",
		PredictedTimes: nil,
		TrainingPaces:  nil,
		VDOT:           0.0,
		DeviceName:     "",
		DeviceUnitID:   "",
		DeviceProdID:   "",
		RunScore:       0.0,
		VO2max:         0.0,
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

	p.LapDist, p.LapTime, p.LapCal, p.LapPace, p.TotalDistance = processFitLap(runLaps, toEnglish)

	// Get the start time.
	p.Titletext += time.Unix(p.TimeStamps[0], 0).Format(time.UnixDate)

	// Calculate the summary string information.
	p.DispTotalDistance, p.TotalPace, p.ElapsedTime, p.TotalCal, p.AvgPower, p.StartDateStamp,
		p.EndDateStamp = createStats(toEnglish, p.TotalDistance, p.TimeStamps, p.LapCal)

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

// Allow the user to shutdown the server from a target in the user interface client.
func stopHandler(w http.ResponseWriter, r *http.Request) {
	// Nothing graceful about this exit.  Just bail out.
	os.Exit(0)
	return
}


func main() {
	desktop.Open("http://localhost:8080")
	// Serve static files if the prefix is "static".
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Handle normal requests.
	http.HandleFunc("/", pageloadHandler)
	http.HandleFunc("/getplot", plotHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/getruns", dbHandler)
	http.HandleFunc("/selectrun", dbSelectHandler)
	//Listen on port 8080
	//fmt.Println("Server starting on port 8080.")
	http.ListenAndServe(":8080", nil)
}
