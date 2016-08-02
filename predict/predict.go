//
// Package predict makes running performance estimates based on VO2max formulas.
//
package predict

import (
	//"fmt"
	"errors"
	"math"
)

// Oxygen Cost Formula:
// O2cost = 0.000104 x (velocity)2 + 0.182258 X (velocity) - 4.60
// The reference to "velocity" is, of course, referring to the running velocity.
// The velocity is expressed in meters per minute (m/min).
// The oxygen cost is expressed in milliliters of oxygen per kilogram of the runner's
// weight per minute (ml/kg/min).
//
// Ref: J. Daniels, R. Fitts and G. Sheehan  The Conditioning for Distance Running--
//  the Scientific Aspects (John Wiley & Sons, New York, 1978)]
func calcO2cost(v float64) (O2cost float64) {
	O2cost = 0.000104*math.Pow(v, 2) + (0.182258 * v) - 4.60
	return O2cost
}

// Intensity (aka "Drop Dead Formula")
// I = 0.2989558e-0.1932605t + 0.1894393e-0.012778t + 0.8
// The "t" is the amount of time in minutes that a human can run at the calculated
// intensity expressed as a percentage of a person's maximum oxygen uptake capacity.
// Ref: Daniels, Gilbert
func calcIntensity(t float64) float64 {
	Intensity := 0.2989558*math.Exp(-0.1932605*t) + 0.1894393*math.Exp(-0.012778*t) + 0.8
	return Intensity
}

// CalcVO2max calculates V02max using the Daniel's and Gilbert formula.
// The velocity is expressed in meters per minute.
// The time is expressed in minutes.
func CalcVO2max(v, t float64) float64 {
	O2Cost := calcO2cost(v)
	Intensity := calcIntensity(t)
	VO2max := O2Cost / Intensity
	return VO2max
}

// Takes a VO2 measurement and converts it to a velocity in m/min.
func vo2ToPace(vo2Val float64) float64 {
	return 29.54 + 5.000663*vo2Val - 0.007546*(vo2Val*vo2Val)
}

// Bisect finds the root of the function (in this case, CalcVO2max) by bisecting
// https://en.wikipedia.org/wiki/Bisection_method
// The function argument (fn) signature is specific to CalcVO2max.
func Bisect(fn func(float64, float64) float64, a float64, b float64, tol float64,
	maxIter int, raceLengthMeters float64) (c float64, err error) {

	n := 1
	for n <= maxIter {
		c := (a + b) / 2.0 // new midpoint
		if fn(c, raceLengthMeters) == 0 || (b-a)/2.0 < tol {
			return c, nil
		} else {
			n++
			if fn(c, raceLengthMeters)*fn(a, raceLengthMeters) > 0.0 {
				a = c
			} else {
				b = c
			}
		}
	}
	return math.NaN(), errors.New("Failed to converge.")
}

// Daniels calculates a predicted race time using the Daniel's Gilbert VO2max criteria.
func Daniels(runLengthMeters float64, elapsedTime float64,
	raceLengthMeters float64) (tOut float64, VO2max float64, err error) {
	// Inputs are:
	// a run length in meters and the elapsed time in minutes
	//
	// Outputs are:
	// tOut represents the number of minutes predicted for raceLengthMeters
	// VO2max is expressed in milliliters of oxygen per kilogram of the runner's
	// weight per minute (ml/kg/min).
	// err will be set if the solver failed to converge.

	// Calculate the runner's VO2max based on a current run/race.
	vRun := runLengthMeters / elapsedTime // run velocity meters/min
	VO2max = CalcVO2max(vRun, elapsedTime)
	// For a race prediction, we need to solve the VO2max equation for time,
	// given a VO2max either measured or from a training run or race (above) and a
	// distance for the race.
	// We'll use a simple bisection root solver method.
	// Implementation: What sorcery is this?
	// Go allows function to be passed as arguments.
	// Ref: https://tour.golang.org/moretypes/24
	fn := func(raceTimeinMinutes, raceLengthMeters float64) float64 {
		v := raceLengthMeters / raceTimeinMinutes
		return CalcVO2max(v, raceTimeinMinutes) - VO2max
	}
	a := 1.0       // Initial low guess for solution e.g. 1 minute 400 m
	b := 300.0     // Initial high guess for solution e.g. 5 hour marathon in minutes.
	tol := 0.01    // Solution tolerance 0.05 min = 3 secs = margin of error.
	maxIter := 100 //Fail if not converged after maxIter loops.
	root, err := Bisect(fn, a, b, tol, maxIter, raceLengthMeters)
	tOut = root
	return tOut, VO2max, err

}

// Races predicts race times using the Daniel's Gilbert VO2max criteria.
func Races(runLengthMeters float64, elapsedTime float64) (PredictedTimes map[string]float64, VDOT float64, err error) {
	// Make predictions.  Tell the future.  :)
	//
	// Inputs are: a run length in meters and elapsed time in minutes.
	//
	// Outputs are:
	// tOut represents an array of the number of minutes predicted for common race lengths
	// VO2max is expressed in milliliters of oxygen per kilogram of the runner's
	// weight per minute (ml/kg/min).
	// err will be set if the solver failed to converge.
	PredictedTimes = make(map[string]float64)
	PredictedTimes["400"], VDOT, _ = Daniels(runLengthMeters, elapsedTime, 400.0)
	PredictedTimes["800"], _, _ = Daniels(runLengthMeters, elapsedTime, 800.0)
	PredictedTimes["1 mi"], _, _ = Daniels(runLengthMeters, elapsedTime, 1609.34)
	PredictedTimes["5k"], _, _ = Daniels(runLengthMeters, elapsedTime, 5000.0)
	PredictedTimes["10k"], _, _ = Daniels(runLengthMeters, elapsedTime, 10000.0)
	PredictedTimes["10mi"], _, _ = Daniels(runLengthMeters, elapsedTime, 16093.4)
	PredictedTimes["Half Marathon"], _, _ = Daniels(runLengthMeters, elapsedTime, 21097.4)
	PredictedTimes["25k"], _, _ = Daniels(runLengthMeters, elapsedTime, 25000.0)
	PredictedTimes["30k"], _, _ = Daniels(runLengthMeters, elapsedTime, 30000.0)
	PredictedTimes["Marathon"], _, _ = Daniels(runLengthMeters, elapsedTime, 42195.0)
	return PredictedTimes, VDOT, err
}

// TrainingPaces calculates training paces based on percentages of VO2max.
func TrainingPaces(VO2max float64) (easyPace, maraPace, thresholdPace, intervalPace, repPace float64) {
	// Calculate training paces in meters/min
	easyPace = vo2ToPace(VO2max * 0.7)       // 70% vo2max
	maraPace = vo2ToPace(VO2max * 0.82)      // 82% vo2max
	thresholdPace = vo2ToPace(VO2max * 0.88) // 88% vo2max
	intervalPace = vo2ToPace(VO2max * 0.98)  // 98% vo2max
	repPace = vo2ToPace(VO2max * 1.05)       // 105% vo2max
	return

}
