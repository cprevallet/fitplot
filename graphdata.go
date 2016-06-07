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

func unitCvt(val float64, valtype string, toEnglish bool) (cvtVal float64) {
	switch {
        case valtype == "distance":
	  if toEnglish {
		cvtVal = val * metersToMiles
	    } else {
		cvtVal = val * metersToKm
	   }
        case valtype == "pace" :
	  if toEnglish {
		cvtVal = val * paceToEnglish
	    } else {
		cvtVal = val * paceToMetric
	   }
        case valtype == "altitude":
	  if toEnglish {
		cvtVal = val * metersToFt
	    } else {
		cvtVal = val
	   }
        case valtype == "cadence":
	  cvtVal = val
        }
  return cvtVal
}

func getDvsA(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
	    x = unitCvt(record.Distance, "distance", toEnglish)
	    y = unitCvt(record.Altitude, "altitude", toEnglish)
	    coordpair := []float64{x,y}
            data = append(data, coordpair)
        }
        return
}

func getDvsC(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
	    x = unitCvt(record.Distance, "distance", toEnglish)
	    y = unitCvt(float64(record.Cadence), "cadence", toEnglish)
	    coordpair := []float64{x,y}
            data = append(data, coordpair)
        }
        return
}

func getDvsP(runRecs []fit.Record, toEnglish bool) (data [][]float64) {
	var x float64
	var y float64
        for _, record := range runRecs {
	    x = unitCvt(record.Distance, "distance", toEnglish)
	    if record.Speed > 0.0 {
	      y = unitCvt(1.0/record.Speed, "pace", toEnglish)
	    } else {
	      if toEnglish {y=20.0} else {y=12.0}
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

