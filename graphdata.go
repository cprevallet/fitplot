package main

//
// Open a fit file and organize as datapoints suitable to plot.
//

import (
  "github.com/jezard/fit"
)

var metersToMiles float64 = 0.00062137119 //meter -> mile
var metersToKm float64 = 0.001  //meter -> km
var metersToFt float64 = 3.2808399 //meter ->ft
var paceToEnglish float64 = 26.8224  // sec/meter -> min/mile
var paceToMetric float64 = 16.666667 // sec/meter -> min/km

var toEnglish bool = true


func getDvsA(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
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

func getDvsC(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
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

func getDvsP(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
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

func getlatlong(runRecs []fit.Record) (data []map[string]float64) {
        for _, record := range runRecs {
	    mapPos := map[string]float64{"lat": record.Position_lat, "lng": record.Position_long}
            data = append(data, mapPos)
        }
        return
}

