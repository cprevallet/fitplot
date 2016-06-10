package main

//
// Open a fit file and organize as datapoints suitable to plot.
//

import (
  "github.com/jezard/fit"
  "github.com/cprevallet/fitplot/stats"
//  "fmt"
  "math"
  "strconv"
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

// Slice up a structure.
func unpackRecs( runRecs []fit.Record) (timestamp []int64, distance []float64, altitude []float64, cadence []float64, speed []float64, lat []float64, lng []float64) {
  // TODO Consider make([]float64, len(record.Distance), len(record.Distance)
  // Should reduce time to allocate on each iteration.
  for _, record := range runRecs {
//	  fmt.Println("Timestamp",time.Unix(record.Timestamp, 0).UTC().Format(time.RFC3339))
	  timestamp = append(timestamp, record.Timestamp)
	  distance = append(distance, record.Distance)
	  altitude = append(altitude, record.Altitude)
	  cadence = append(cadence, float64(record.Cadence))
	  speed = append(speed, record.Speed)
	  lat = append(lat, record.Position_lat)
	  lng = append(lng, record.Position_long)
	}
	return	 
  }

// Take two slices and pair the individual elements into x-y points.
func createPlotCoordinates(xSlice []float64, ySlice []float64)(data [][]float64) {
  for i, _ := range xSlice {
    coordpair := []float64{xSlice[i], ySlice[i]}
    data = append(data, coordpair)
  }
  return
}

// Convert two arrays into a map used by Google maps.
func getMapCoordinates(latSlice []float64, lngSlice []float64) (data []map[string]float64) {
  for i, _ := range latSlice {
    mapPos := map[string]float64{"lat": latSlice[i], "lng": lngSlice[i]}
    data = append(data, mapPos)
  }
  return
}

// Main entry point
// Convert the record structure to slices and maps suitable for use in the user interface.
func processFitRecord(runRecs []fit.Record, toEnglish bool)( mapData []map[string]float64,Y0Pairs [][]float64, Y1Pairs [][]float64, Y2Pairs [][]float64, dispTimestamp[]int64 ) {

    // Get slices from the runRecs structure.
    timestamp, distance, altitude, cadence, speed, lat, lng := unpackRecs(runRecs)
    // Speed -> pace
    var pace []float64
    for i, _ := range(speed) {
      if speed[i] > 1.8 {  //m/s
	pace = append(pace, 1.0/speed[i])
      } else {
	pace = append(pace, 0.56 )  // s/m = 15 min/mi
      }
    }
    
    // Clean up the data statistically before displaying.
    outIdxs := markOutliers(pace)
    for _, item := range markOutliers(lat) {
      outIdxs = append(outIdxs, item)
    }
    for _, item := range markOutliers(lng) {
      outIdxs = append(outIdxs, item)
    }
    timestamp_clean := removeOutliersInt(timestamp, outIdxs)
    distance_clean := removeOutliers(distance, outIdxs)
    pace_clean := removeOutliers(pace, outIdxs)
    altitude_clean := removeOutliers(altitude, outIdxs)
    cadence_clean := removeOutliers(cadence, outIdxs)
    lat_clean := removeOutliers(lat, outIdxs)
    lng_clean := removeOutliers(lng, outIdxs)  
    
    // Convert the units for the slices that have them.
    dispDistance := convertUnits(distance_clean, "distance", toEnglish)
    dispPace := convertUnits(pace_clean, "pace", toEnglish)
    dispAltitude := convertUnits(altitude_clean, "altitude", toEnglish)
    dispCadence := convertUnits(cadence_clean, "cadence", toEnglish)
    dispTimestamp = timestamp_clean
    
    //Return the values used in the user interface.
    Y0Pairs = createPlotCoordinates(dispDistance, dispPace)
    Y1Pairs = createPlotCoordinates(dispDistance, dispAltitude)
    Y2Pairs = createPlotCoordinates(dispDistance, dispCadence)
    mapData = getMapCoordinates(lat_clean, lng_clean)
    
    return
  }

func markOutliers( x []float64 ) (outliersIdx []int) {
   // Create a list of indexs where the value of x is outside of the 
   // 99.7% (3 sigma) expected value assuming a normal distribution of x.
   // In English, find the "unusual" points.
   mean := stats.Sum(x) / float64(len(x))
   sigma := stats.StdDev(x, mean)
   upperLimit := mean + (3.0 * sigma)
   lowerLimit := mean - (3.0 * sigma)
   for i, _ := range(x) {
     if x[i] < lowerLimit || x[i] > upperLimit {
         outliersIdx = append(outliersIdx, i)
     }
  }
  return outliersIdx
}

func removeOutliers(x[]float64, outliersIdx []int) (z[]float64)  {
  // Remove values in x if it's index matches one in the list of outliers.
  for i, item := range(x) {
    found := false
    for _, idx := range(outliersIdx) {
      if i == idx { found = true }
      }
    if !found {z = append(z, item)}
    }
  return z
}

func removeOutliersInt(x[]int64, outliersIdx []int) (z[]int64)  {
  // Remove values in x if it's index matches one in the list of outliers.
  for i, item := range(x) {
    found := false
    for _, idx := range(outliersIdx) {
      if i == idx { found = true }
      }
    if !found {z = append(z, item)}
    }
  return z
}


// Main entry point
// Convert the record structure to slices and maps suitable for use in the user interface.
func processFitLap(runLaps []fit.Lap, toEnglish bool) (LapDist []float64, LapTime []string, LapCal []float64, LapPace []string){
  for _, item := range(runLaps) {
    dist := unitCvt(item.Total_distance, "distance", toEnglish)
    cal := float64(item.Total_calories)
    // Seconds to "min:sec"
    laptime_str := decimalTimetoMinSec(float64(item.Total_elapsed_time/60.0))
    // Calculate pace string.
    pace := item.Total_elapsed_time/60.0/dist
    //pace = unitCvt(pace, "pace", toEnglish)
    pace_str := decimalTimetoMinSec(pace)
    LapDist = append(LapDist, dist)
    LapCal = append(LapCal, cal)
    LapPace = append(LapPace, pace_str)
    LapTime = append(LapTime, laptime_str)
   }
   return LapDist, LapTime, LapCal, LapPace
}

// Convert decimal minutes to m:ss.
func decimalTimetoMinSec(in float64) (out string) {
  in_min := int(math.Floor(in))
  in_sec := int((in - float64(in_min))* 60)
  out = strconv.Itoa(in_min) + ":"
  if in_sec < 10 {
    out = out + "0" + strconv.Itoa(in_sec)
  } else {
    out = out + strconv.Itoa(in_sec)
  }
  return out
}
