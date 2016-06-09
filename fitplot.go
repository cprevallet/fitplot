package main

import (
	"encoding/json"
	"fmt"
        "github.com/cprevallet/fitplot/tcx"
        "github.com/jezard/fit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
//	"net/http/httputil"
	"time"
)

var uploadFname string = ""

//Compile templates on start for better performance.
var templates = template.Must(template.ParseFiles("tmpl/fitplot.html"))

//Display the named template.
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

// Handle requests to "/".
func pageloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// Load page.
        	fmt.Println("pageloadHandler Received Request")
		display(w, "fitplot", nil)
	}
	if r.Method == "POST"  {
		// File load request
        	fmt.Println("pageloadHandler POST Received Request")
		uploadHandler(w, r)
		display(w, "fitplot", nil)
	}

}

//Upload a copy the fit file to a temporary local directory.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["myfiles"]
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
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
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("uploadHandler Received Request")
}

func plotHandler(w http.ResponseWriter, r *http.Request) {

	type Plotvals struct {
	    Titletext string
	    XName string
	    Y0Name string
	    Y1Name string
	    Y2Name string 
	    Y0coordinates [][]float64
	    Y1coordinates [][]float64
	    Y2coordinates [][]float64
            Latlongs[] map[string]float64 
            LapDist []float64
            LapTime []string
            LapCal  []float64
            LapPace []string
	}
	var timeStamps []int64
        var xStr string = "Distance "
        var y0Str string = "Pace "
        var y1Str string = "Elevation"
        var y2Str string = "Cadence "
	var runRecs []fit.Record
	var runLaps []fit.Lap

	//what has the user selected for unit system?

	toEnglish = true
	param1s := r.URL.Query()["toEnglish"];
	if param1s != nil {
		if (param1s[0] == "true")  {toEnglish = true}
		if (param1s[0] == "false") {toEnglish = false}
	}

//	dump, err := httputil.DumpRequest(r, true)
//		if err != nil {
//			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
//			return
//		}

//		fmt.Printf("%s\n\n", dump)

        //Read file. uploadFname gets set in uploadHandler.
	b, _ := ioutil.ReadFile(uploadFname)
	rslt := http.DetectContentType(b)
	switch {
        case rslt == "application/octet-stream":
	    // filetype is FIT, or at least it could be?
	    fitStruct := fit.Parse(uploadFname, false)
	    runRecs = fitStruct.Records
	    runLaps = fitStruct.Laps
 
	    //for _, lap := range fitStruct.Laps {
	    //  fmt.Printf("%+v\n", lap)
	    //}
        case rslt == "text/xml; charset=utf-8" :
	    // filetype is TCX or at least it could be?
	    db, _ := tcx.ReadTCXFile(uploadFname)
	    //We cleverly convert the values of interest into a structure we already 
	    //can handle.
	    runRecs = cvtToFitRecs(db)
	    runLaps = cvtToFitLaps(db)
	    //Here's where the heavy lifting of pulling tracks and performance information
            //from (portions of) the fit file into something we can view is done.
	    //for i, _ := range db.Acts.Act {
	//	for _, lap := range db.Acts.Act[i].Laps {
	//	fmt.Printf("%+v\n", lap)
	//	}
	  
	}
	fmt.Printf("%+v\n", runLaps)
        //Build the variable strings based on unit system.
        if toEnglish {
            xStr = xStr + "(mi)"
            y0Str = y0Str + "(min/mi)"
            y1Str = y1Str + "(ft)"
            y2Str = y2Str + "(bpm)"
        } else {
            xStr = xStr + "(km)"
            y0Str = y0Str + "(min/km)"
            y1Str = y1Str + "(m)"
            y2Str = y2Str + "(bpm)"
        }

        //Create an object to contain various plot values.
        p := Plotvals {Titletext: "", 
                XName: xStr, 
                Y0Name: y0Str,
                Y1Name: y1Str,
                Y2Name: y2Str,
                Y0coordinates: nil,
                Y1coordinates: nil,
                Y2coordinates: nil,
                Latlongs: nil,
		LapDist: nil,
		LapTime: nil,
		LapCal: nil,
		LapPace: nil,
        }


	//Here's where the heavy lifting of pulling tracks and performance information
        //from (portions of) the fit file into something we can view is done.
	p.Latlongs ,p.Y0coordinates, p.Y1coordinates, p.Y2coordinates, timeStamps =
	  processFitRecord(runRecs, toEnglish)
	
	p.LapDist, p.LapTime, p.LapCal, p.LapPace = processFitLap(runLaps, toEnglish)
	
	//Get start time.
	p.Titletext += time.Unix(timeStamps[0], 0).Format(time.UnixDate)
	
        //Convert to json.
        js, err := json.Marshal(p)

        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }


        fmt.Println("plotHandler Received Request")
        w.Header().Set("Content-Type", "text/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //Send
        w.Write(js)
}

func main() {
	http.HandleFunc("/", pageloadHandler)    
	http.HandleFunc("/getplot", plotHandler) 
	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
