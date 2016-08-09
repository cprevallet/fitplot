//
// Fitplot provides a webserver used to process .fit and .tcx files.
//
package main

import (
	"encoding/json"
	"fmt"
	"github.com/cprevallet/fitplot/desktop"
	"github.com/cprevallet/fitplot/persist"
	"github.com/cprevallet/fitplot/tcx"
	"github.com/jezard/fit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
var uploadFname = ""

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
	// Parse the multipart form in the request.
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a ref to the parsed multipart form.
	m := r.MultipartForm

	// Get the *fileheaders.
	_ = "breakpoint"
	myfile := m.File["file"]
	// Get a handle to the actual file.
	file, err := myfile[0].Open()
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dst, err := ioutil.TempFile("", "example")
	uploadFname = ""
	uploadFname = dst.Name()
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Copy the uploaded file to the destination file.
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println("uploadHandler Received Request")
	// Do some low-level stuff to retrieve the file name and byte array.
	_ = "breakpoint"
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fName := handler.Filename
	// Persist it to the database.
	dbHandler(w, r, fName, fBytes)

}

// Initialize the database used to store run files if one doesn't exist.
func dbHandler(w http.ResponseWriter, r *http.Request, uploadFname string, file[]byte) {
	db, _ := persist.ConnectDatabase("fitplot", "./")
	persist.InsertNewRecord(db, uploadFname, file)
	db.Close()
}

//
// Return information about the runtime environment.
//
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
		TotalDistance  string
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
	var runRecs []fit.Record
	var runLaps []fit.Lap
	var c0Str, c1Str, c2Str, c3Str, c4Str string
	var fitStruct fit.FitFile
	var db *tcx.TCXDB

	// User hasn't uploaded a file yet?  Avoid a panic.
	if uploadFname == "" {
		http.Error(w, "No file loaded.", http.StatusConflict)
		return
	}
	// What has the user selected for unit system?
	toEnglish = true
	param1s := r.URL.Query()["toEnglish"]
	if param1s != nil {
		if param1s[0] == "true" {
			toEnglish = true
		}
		if param1s[0] == "false" {
			toEnglish = false
		}
	}
	// What race time/distance has the user entered?
	var racedist float64
	var racehours, racemins, racesecs int64
	param2s := r.URL.Query()["racedist"]
	if param2s != nil {
		racedist, _ = strconv.ParseFloat(param2s[0], 64)
	} else {
		racedist = 5000.0
	}
	param3s := r.URL.Query()["racehours"]
	if param3s != nil {
		racehours, _ = strconv.ParseInt(param3s[0], 10, 64)
	} else {
		racehours = 0
	}
	param4s := r.URL.Query()["racemins"]
	if param4s != nil {
		racemins, _ = strconv.ParseInt(param4s[0], 10, 64)
	} else {
		racemins = 25
	}
	param5s := r.URL.Query()["racesecs"]
	if param5s != nil {
		racesecs, _ = strconv.ParseInt(param5s[0], 10, 64)
	} else {
		racesecs = 0
	}

	// Calculate analysis results based on segment or whole run?
	useSegment := false
	param6s := r.URL.Query()["useSegment"]
	if param6s != nil {
		if param6s[0] == "true" {
			useSegment = true
		}
		if param6s[0] == "false" {
			useSegment = false
		}
	}
	// What split time/distance has the user entered?
	var splitdist float64
	var splithours, splitmins, splitsecs int64
	param7s := r.URL.Query()["splitdist"]
	if param7s != nil {
		splitdist, _ = strconv.ParseFloat(param7s[0], 64)
	} else {
		splitdist = 5000.0
	}
	param8s := r.URL.Query()["splithours"]
	if param8s != nil {
		splithours, _ = strconv.ParseInt(param8s[0], 10, 64)
	} else {
		splithours = 0
	}
	param9s := r.URL.Query()["splitmins"]
	if param9s != nil {
		splitmins, _ = strconv.ParseInt(param9s[0], 10, 64)
	} else {
		splitmins = 25
	}
	param10s := r.URL.Query()["splitsecs"]
	if param10s != nil {
		splitsecs, _ = strconv.ParseInt(param10s[0], 10, 64)
	} else {
		splitsecs = 0
	}
	/*
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%s\n\n", dump)
	*/
	// Read file. uploadFname gets set in uploadHandler.
	b, _ := ioutil.ReadFile(uploadFname)
	rslt := http.DetectContentType(b)
	switch {
	case rslt == "application/octet-stream":
		// Filetype is FIT, or at least it could be?
		fitStruct = fit.Parse(uploadFname, false)
		runRecs = fitStruct.Records
		runLaps = fitStruct.Laps

	case rslt == "text/xml; charset=utf-8":
		// Filetype is TCX or at least it could be?
		db, _ = tcx.ReadTCXFile(uploadFname)
		//		if err != nil {
		//			fmt.Printf("Error parsing file", err)
		//		}

		// We cleverly convert the values of interest into a structures we already
		// can handle.
		runRecs = tcx.CvtToFitRecs(db)
		runLaps = tcx.CvtToFitLaps(db)
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
		TotalDistance:  "",
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

	if rslt == "application/octet-stream" {
		p.DeviceName = fitStruct.DeviceInfo[0].Manufacturer
		p.DeviceProdID = fitStruct.DeviceInfo[0].Product
		p.DeviceUnitID = fmt.Sprint(fitStruct.DeviceInfo[0].Serial_number)
	}
	if rslt == "text/xml; charset=utf-8" {
		p.DeviceName, p.DeviceUnitID, p.DeviceProdID = tcx.DeviceInfo(db)
	}

	// Here's where the heavy lifting of pulling tracks and performance information
	// from (portions of) the fit file into something we can view is done.
	p.Latlongs, p.TimeStamps, p.DispDistance, p.DispPace, p.DispAltitude, p.DispCadence =
		processFitRecord(runRecs, toEnglish)

	p.LapDist, p.LapTime, p.LapCal, p.LapPace = processFitLap(runLaps, toEnglish)

	// Get the start time.
	p.Titletext += time.Unix(p.TimeStamps[0], 0).Format(time.UnixDate)

	// Calculate the summary string information.
	p.TotalDistance, p.TotalPace, p.ElapsedTime, p.TotalCal, p.AvgPower, p.StartDateStamp,
		p.EndDateStamp = createStats(toEnglish, p.DispDistance, p.TimeStamps, p.LapCal)

	// Calculate the analysis page.
	p.PredictedTimes, p.VDOT, p.VO2max, p.RunScore, p.TrainingPaces = createAnalysis(toEnglish,
		useSegment, p.DispDistance, p.TimeStamps, splitdist, splithours, splitmins,
		splitsecs, racedist, racehours, racemins, racesecs)

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
	//Listen on port 8080
	//fmt.Println("Server starting on port 8080.")
	http.ListenAndServe(":8080", nil)
}
