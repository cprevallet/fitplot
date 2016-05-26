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

var metersToMiles float64 = 0.00062137119 //meter -> mile
var metersToKm float64 = 0.001  //meter -> km
var metersToFt float64 = 3.2808399 //meter ->ft
var paceToEnglish float64 = 26.8224  // sec/meter -> min/mile
var paceToMetric float64 = 16.666667 // sec/meter -> min/km

var toEnglish bool = true

//
//  Main
//
func main() {
	http.HandleFunc("/",foo)
	http.ListenAndServe(":3000", nil)
	}

func getDvsA(fitStruct fit.FitFile, toEnglish bool) (data [][]float64) {

	var x float64
	var y float64
        for _, record := range fitStruct.Records {
	    if toEnglish {
		x = record.Distance * metersToMiles
                y = record.Altitude * metersToFt
	    } else {
		x = record.Distance * metersToKm
                y = record.Altitude
	   }
	    coordpair := []float64{x,y}
            data = append(data, coordpair)
        }
        return
}

func getDvsC(fitStruct fit.FitFile, toEnglish bool) (data [][]float64) {

	var x float64
	var y float64
        for _, record := range fitStruct.Records {
	    if toEnglish {
		x = record.Distance * metersToMiles
                y = float64(record.Cadence)
	    } else {
		x = record.Distance * metersToKm
                y = float64(record.Cadence)
	   }
	    coordpair := []float64{x,y}
            data = append(data, coordpair)
        }
        return
}

func getDvsP(fitStruct fit.FitFile, toEnglish bool) (data [][]float64) {

	var x float64
	var y float64
        for _, record := range fitStruct.Records {
	    if toEnglish {
		x = record.Distance * metersToMiles
                if record.Speed > 0.0 {
                    y = 1.0/record.Speed * paceToEnglish
                } else {
                    y = 20.0  //min/mile...walking pace if standing still...gotta choose something
                }
	    } else {
		x = record.Distance * metersToKm
                if record.Speed > 0.0 {
                    y = 1.0/record.Speed * paceToMetric
                } else {
                    y = 12.0  //min/km...walking pace if standing still...gotta choose something
                }
	   }
	    coordpair := []float64{x,y}
            data = append(data, coordpair)
        }
        return
}



func foo(w http.ResponseWriter, r *http.Request) {
        
	//Read .fit file.
	//var fitFname = "/home/penguin/work/src/github.com/csp/fitplot/65L84100.FIT"
        var fitFname = "/home/penguin/work/src/github.com/csp/fitplot/64IA2815.FIT"
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
