package main

import (
  "fmt"
  "github.com/jezard/fit"
  "encoding/json"
  "net/http"
)

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

//
//  Main
//
func main() {
	http.HandleFunc("/",foo)
	http.ListenAndServe(":3000", nil)
	}

func getDvsA(fitStruct fit.FitFile) (data [][]float64) {

        for _, record := range fitStruct.Records {
	    coordpair := []float64{record.Distance,float64(record.Altitude)}
            data = append(data, coordpair)
        }
        return
}

func getDvsC(fitStruct fit.FitFile) (data [][]float64) {

        for _, record := range fitStruct.Records {
	    coordpair := []float64{record.Distance,float64(record.Cadence)}
            data = append(data, coordpair)
        }
        return
}

func getDvsS(fitStruct fit.FitFile) (data [][]float64) {

        for _, record := range fitStruct.Records {
	    coordpair := []float64{record.Distance, float64(record.Speed)}
            data = append(data, coordpair)
        }
        return
}



func foo(w http.ResponseWriter, r *http.Request) {
        
	//Read .fit file
	//var fitFname = "/home/penguin/work/src/github.com/csp/fitplot/65L84100.FIT"
        var fitFname = "/home/penguin/work/src/github.com/csp/fitplot/64IA2815.FIT"
        var fitStruct fit.FitFile
        fitStruct = fit.Parse(fitFname, false)

        //Convert to a form (x-y pairs) for graph
	spddata := getDvsS(fitStruct)
	altdata := getDvsA(fitStruct)
	caddata := getDvsC(fitStruct)

	//Create an object to contain various plot values
	p := Plotvals {Titletext: "Distance Graph", 
		XName: "Distance (m)", 
		Y0Name: "Speed (m/s)", 
		Y1Name: "Altitude (m)", 
		Y2Name: "Cadence (bpm)", 
		Y0coordinates: nil,
		Y1coordinates: nil,
		Y2coordinates: nil,
	}
	p.Y0coordinates = spddata
	p.Y1coordinates = altdata
	p.Y2coordinates = caddata

	//Convert to json
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
