package main

import (
	"encoding/json"
	"fmt"
        "github.com/jezard/fit"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

var fitFname string = ""


//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/fitplot.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//Display the unitialized graph.  After the user hits the load button,
//copy the fit file to a temporary local directory.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the unitialized graph.
	case "GET":
		display(w, "fitplot", nil)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
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
			fitFname = "" 
			fitFname = dst.Name()
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
		//display success message.
		display(w, "fitplot", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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
	}

        //Read .fit file.
        var fitStruct fit.FitFile
        fitStruct = fit.Parse(fitFname, false)

        //Build the variable strings based on unit system.
        var xStr string = "Distance "
        var y0Str string = "Pace "
        var y1Str string = "Altitude "
        var y2Str string = "Cadence "
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
        p := Plotvals {Titletext: "Distance Graph", 
                XName: xStr, 
                Y0Name: y0Str,
                Y1Name: y1Str,
                Y2Name: y2Str,
                Y0coordinates: nil,
                Y1coordinates: nil,
                Y2coordinates: nil,
        }

        //Convert to a form (x-y pairs) for graph.
        p.Y0coordinates = getDvsP(fitStruct, toEnglish)
        p.Y1coordinates = getDvsA(fitStruct, toEnglish)
        p.Y2coordinates = getDvsC(fitStruct, toEnglish)

        //Convert to json.
        js, err := json.Marshal(p)

        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        fmt.Println("Received Request")
        w.Header().Set("Content-Type", "text/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //Send
        w.Write(js)
        }

func mapHandler(w http.ResponseWriter, r *http.Request) {

        //Read .fit file.
        var fitStruct fit.FitFile
        fitStruct = fit.Parse(fitFname, false)

        //Convert to a latitude longitude for graph.
        var coords[] map[string]float64 
	coords = getlatlong(fitStruct)

        //Convert to json.
        js, err := json.Marshal(coords)

        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("What the hell?")
                return
        }
        fmt.Println("Received Request")
        w.Header().Set("Content-Type", "text/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //Send
        w.Write(js)
        }

func main() {
	http.HandleFunc("/upload", uploadHandler) //url associated with UI
	http.HandleFunc("/getplot", plotHandler) //url for server to supply the plot data
	http.HandleFunc("/getmap", mapHandler) //url for server to supply the map data
	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
