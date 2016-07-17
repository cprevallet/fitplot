package main

//
// Open a fit file and organize as datapoints suitable to plot.
//

import (
	"github.com/cprevallet/fitplot/predict"
	"github.com/cprevallet/fitplot/stats"
	"github.com/cprevallet/fitplot/strutil"
	"github.com/jezard/fit"
	//"fmt"
	"math"
	"strconv"
	"time"
)

// Conversions
var metersToMiles float64 = 0.00062137119 // meter -> mile
var metersToKm float64 = 0.001            // meter -> km
var metersToFt float64 = 3.2808399        // meter ->ft
var paceToEnglish float64 = 26.8224       // sec/meter -> min/mile
var paceToMetric float64 = 16.666667      // sec/meter -> min/km
var stridestoSteps float64 = 2.0          // strides/min -> steps/min (for bipeds)
// Unit system set by user.
var toEnglish bool = true

// Do unit conversion for slices of type valtype.
func convertUnits(vals []float64, valtype string, toEnglish bool) (cvtVals []float64) {
	for _, val := range vals {
		converted := unitCvt(val, valtype, toEnglish)
		cvtVals = append(cvtVals, converted)
	}
	return cvtVals
}

// Convert a single value (e.g. from a slice).
func unitCvt(val float64, valtype string, toEnglish bool) (cvtVal float64) {
	switch {
	case valtype == "distance":
		if toEnglish {
			cvtVal = val * metersToMiles
		} else {
			cvtVal = val * metersToKm
		}
	case valtype == "pace":
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
		cvtVal = val * stridestoSteps
	}
	return cvtVal
}

// Slice up a structure.
func unpackRecs(runRecs []fit.Record) (timestamp []int64, distance []float64,
	altitude []float64, cadence []float64, speed []float64,
	lat []float64, lng []float64) {
	for _, record := range runRecs {
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

// Convert two arrays into a map used by Google maps.
func getMapCoordinates(latSlice []float64, lngSlice []float64) (data []map[string]float64) {
	for i, _ := range latSlice {
		mapPos := map[string]float64{"lat": latSlice[i], "lng": lngSlice[i]}
		data = append(data, mapPos)
	}
	return
}

// This is a main entry point
// Convert the record structure to slices and maps suitable for use in the user interface.
func processFitRecord(runRecs []fit.Record, toEnglish bool) (mapData []map[string]float64, dispTimestamp []int64, dispDistance []float64, dispPace []float64, dispAltitude []float64, dispCadence []float64) {
	// Get slices from the runRecs structure.
	timestamp, distance, altitude, cadence, speed, lat, lng := unpackRecs(runRecs)
	// Speed -> pace
	var pace []float64
	for i, _ := range speed {
		if speed[i] > 1.8 { //m/s
			pace = append(pace, 1.0/speed[i])
		} else {
			pace = append(pace, 0.56) // s/m = 15 min/mi
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
	dispDistance = convertUnits(distance_clean, "distance", toEnglish)
	dispPace = convertUnits(pace_clean, "pace", toEnglish)
	dispAltitude = convertUnits(altitude_clean, "altitude", toEnglish)
	dispCadence = convertUnits(cadence_clean, "cadence", toEnglish)
	dispTimestamp = timestamp_clean

	//Return the values used in the user interface.
	mapData = getMapCoordinates(lat_clean, lng_clean)

	return
}

// Create a list of indexs where the value of x is outside of the
// 99.7% (3 sigma) expected value assuming a normal distribution of x.
// In English, find the "unusual" points.
func markOutliers(x []float64) (outliersIdx []int) {
	mean := stats.Sum(x) / float64(len(x))
	sigma := stats.StdDev(x, mean)
	upperLimit := mean + (3.0 * sigma)
	lowerLimit := mean - (3.0 * sigma)
	for i, _ := range x {
		if x[i] < lowerLimit || x[i] > upperLimit {
			outliersIdx = append(outliersIdx, i)
		}
	}
	return outliersIdx
}

// Remove values in x if it's index matches one in the list of outliers.
func removeOutliers(x []float64, outliersIdx []int) (z []float64) {
	for i, item := range x {
		found := false
		for _, idx := range outliersIdx {
			if i == idx {
				found = true
			}
		}
		if !found {
			z = append(z, item)
		}
	}
	return z
}

// Same as above except for integer x.
func removeOutliersInt(x []int64, outliersIdx []int) (z []int64) {
	for i, item := range x {
		found := false
		for _, idx := range outliersIdx {
			if i == idx {
				found = true
			}
		}
		if !found {
			z = append(z, item)
		}
	}
	return z
}

// This is a main entry point.
// Convert the record structure to slices and maps suitable for use in the user interface.
func processFitLap(runLaps []fit.Lap, toEnglish bool) (LapDist []float64, LapTime []string, LapCal []float64, LapPace []string) {
	for _, item := range runLaps {
		dist := unitCvt(item.Total_distance, "distance", toEnglish)
		cal := float64(item.Total_calories)
		// Seconds to "min:sec"
		laptime_str := strutil.DecimalTimetoMinSec(float64(item.Total_elapsed_time / 60.0))
		// Calculate pace string.
		pace := item.Total_elapsed_time / 60.0 / dist
		//pace = unitCvt(pace, "pace", toEnglish)
		pace_str := strutil.DecimalTimetoMinSec(pace)
		LapDist = append(LapDist, dist)
		LapCal = append(LapCal, cal)
		LapPace = append(LapPace, pace_str)
		LapTime = append(LapTime, laptime_str)
	}
	return LapDist, LapTime, LapCal, LapPace
}

// Create the summary statistics strings.
func createStats(toEnglish bool, DispDistance []float64, TimeStamps []int64,
	LapCal []float64) (totalDistance string, totalPace string,
	elapsedTime string, totalCal string, avgPower string, startDateStamp string,
	endDateStamp string) {

	// Calculate run start and end times.
	startDateStamp = time.Unix(TimeStamps[0], 0).Format(time.UnixDate)
	endDateStamp = time.Unix(TimeStamps[len(TimeStamps)-1], 0).Format(time.UnixDate)
	// Calculate overall distance.
	totalDistance = strconv.FormatFloat(DispDistance[len(DispDistance)-1], 'f', 2, 64)
	if toEnglish {
		totalDistance += " mi"
	} else {
		totalDistance += " km"
	}
	// Calculate mm:ss for totalPace.
	timeDiffinMinutes := (float64(TimeStamps[len(TimeStamps)-1]) - float64(TimeStamps[0])) / 60.0
	decimalPace := float64(timeDiffinMinutes) / DispDistance[len(DispDistance)-1]
	totalPace = strutil.DecimalTimetoMinSec(decimalPace)
	if toEnglish {
		totalPace += " min/mi"
	} else {
		totalPace += " min/km"
	}
	// Calculate hh:mm:ss for elapsedTime.
	elapsedTime = strutil.DecimalTimetoHourMinSec(float64(timeDiffinMinutes))
	// Sum up the lap calories
	totcal := 0.0
	for _, calorie := range LapCal {
		totcal += calorie
	}
	totalCal = strconv.Itoa(int((math.Floor(totcal)))) + " kcal"
	
	// Calculate power expended based on Garmin calculated calories
	power := totcal * 4186.8 / (timeDiffinMinutes * 60.0)
	avgPower = strconv.FormatFloat(power, 'f', 2, 64) + " Watts"
	return
}

// Do a prediction based on this run.
func createAnalysis(toEnglish bool, 
					   useSegment bool, 
					   DispDistance []float64,
					   TimeStamps []int64, 
					   splitdist float64, 
					   splithours, splitmins, splitsecs int64,
					   racedist float64, 
					   racehours, racemins,racesecs int64,				   
  					) (PredictedRaceTimes map[string]string, 
					   VDOT float64,
					   VO2Max float64,
					   RunScore float64,
					   TrainingPaces map[string]string ) {
	// Need to assign variables based on whether the user has selected the entire
	// run or just a run segment (e.g. a split time and distance) to basis the 
	// prediction on.
	// Do the type conversions to get the inputs into forms expected by PredictRaces.
	var d, dist, tStart, tEnd, elapsedTime float64
	if useSegment != true {
		d = DispDistance[len(DispDistance)-1]
		// Need distance back in meters as PredictRaces demands metric units.
		if toEnglish {
			dist = d / metersToMiles
		} else {
			dist = d / metersToKm		
		}
		tStart = float64(TimeStamps[0])
		tEnd = float64(TimeStamps[len(TimeStamps)-1])
		elapsedTime = (tEnd - tStart) / 60.0 
	} else {
		dist = splitdist
		elapsedTime = (float64(splithours) * 60.0) + float64(splitmins) + (float64(splitsecs) / 60.0)
	}
	
	// Calculate the equivalent race times for this run (segment or complete).
	PredictedTimes, v, _ := predict.PredictRaces(dist, elapsedTime)
	VDOT = v

	// Calculate the % of VO2max for this run relative to the provided race information.
	elapsedTimeRace := (float64(racehours) * 60.0) + float64(racemins) + (float64(racesecs) / 60.0)
	velocityRace := racedist / elapsedTimeRace
	VO2Max = predict.CalcVO2max(velocityRace, elapsedTimeRace)	
	RunScore = VDOT / VO2Max * 100.0
	
	// Calculate the training paces.
	easyPace, maraPace, thresholdPace, intervalPace, repPace := predict.TrainingPaces(VO2Max)
	TrainingPaces = make(map[string]string)
	if toEnglish {
		easyPace = 1609.34 / easyPace
		maraPace = 1609.34 / maraPace
		thresholdPace = 1609.34 / thresholdPace
		intervalPace = 1609.34 / intervalPace
		repPace = 1609.34 / repPace
	} else {
		easyPace = 1000.0 / easyPace
		maraPace = 1000.0 / maraPace
		thresholdPace = 1000.0 / thresholdPace
		intervalPace = 1000.0 / intervalPace
		repPace = 1000.0 / repPace
	}
	TrainingPaces["Easy"] = strutil.DecimalTimetoMinSec(easyPace)
	TrainingPaces["Marathon"] = strutil.DecimalTimetoMinSec(maraPace)
	TrainingPaces["Threshold"] = strutil.DecimalTimetoMinSec(thresholdPace)
	TrainingPaces["Interval"] = strutil.DecimalTimetoMinSec(intervalPace)
	TrainingPaces["Repeats"] = strutil.DecimalTimetoMinSec(repPace)
	
	// Convert the times from decimal minutes to hh:mm:ss for the user.
	PredictedRaceTimes = make(map[string]string)
	for key, val := range PredictedTimes {
		PredictedRaceTimes[key] = strutil.DecimalTimetoHourMinSec(val)
	}
	return
}
