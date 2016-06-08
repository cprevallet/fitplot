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

// Do unit conversion for slices of type valtype
func convertUnits(vals []float64, valtype string, toEnglish bool) (cvtVals []float64) {
        for _, val := range vals {
	  converted := unitCvt(val, valtype, toEnglish)
	  cvtVals = append(cvtVals, converted)
	}
	return cvtVals
}

// Convert a single value (e.g. from a slice)
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

//Slice up a structure.
func unpackRecs( runRecs []fit.Record) (distance []float64, altitude []float64, cadence []float64, speed []float64, lat []float64, lng []float64) {
        for _, record := range runRecs {
	  distance = append(distance, record.Distance)
	  altitude = append(altitude, record.Altitude)
	  cadence = append(cadence, float64(record.Cadence))
	  speed = append(speed, record.Speed)
	  lat = append(lat, record.Position_lat)
	  lng = append(lng, record.Position_long)
	}
	return	 
  }

//Take two slices and pair the individual elements into x-y points.
func createPlotCoordinates(xSlice []float64, ySlice []float64)(data [][]float64) {
  for i, _ := range xSlice {
    coordpair := []float64{xSlice[i], ySlice[i]}
    data = append(data, coordpair)
  }
  return
}

//Convert two arrays into a map used by Google maps.
func getMapCoordinates(latSlice []float64, lngSlice []float64) (data []map[string]float64) {
  for i, _ := range latSlice {
    mapPos := map[string]float64{"lat": latSlice[i], "lng": lngSlice[i]}
    data = append(data, mapPos)
  }
  return
}

//Main entry point
func processFitRecord(runRecs []fit.Record, toEnglish bool)( mapData []map[string]float64,Y0Pairs [][]float64, Y1Pairs [][]float64, Y2Pairs [][]float64 ) {

    // Get slices from the runRecs structure.
    distance, altitude, cadence, speed, lat, lng := unpackRecs(runRecs)
    // Speed -> pace
    var pace []float64
    for i, _ := range(speed) {
      if speed[i] > 1.8 {  //m/s
	pace = append(pace, 1.0/speed[i])
      } else {
	pace = append(pace, 0.56 )  // s/m = 15 min/mi
      }
    }
    // Convert the units for the slices that have them.
    dispDistance := convertUnits(distance, "distance", toEnglish)
    dispPace := convertUnits(pace, "pace", toEnglish)
    dispAltitude := convertUnits(altitude, "altitude", toEnglish)
    dispCadence := convertUnits(cadence, "cadence", toEnglish)
    
    //Return the values used in the user interface.
    Y0Pairs = createPlotCoordinates(dispDistance, dispPace)
    Y1Pairs = createPlotCoordinates(dispDistance, dispAltitude)
    Y2Pairs = createPlotCoordinates(dispDistance, dispCadence)
    mapData = getMapCoordinates(lat, lng)
    
    return
  }
