package predict

//
// Make race performance predictions.
//

import (
//  "fmt"
  "math"
  "errors"
)

func calcO2cost (v float64)(O2cost float64) {
// Oxygen Cost Formula:
// O2cost = 0.000104 x (velocity)2 + 0.182258 X (velocity) - 4.60
// The reference to "velocity" is, of course, referring to the running velocity. 
// The velocity is expressed in meters per minute (m/min). 
// The oxygen cost is expressed in milliliters of oxygen per kilogram of the runner's 
// weight per minute (ml/kg/min).
// 
// Ref: J. Daniels, R. Fitts and G. Sheehan  The Conditioning for Distance Running--
//  the Scientific Aspects (John Wiley & Sons, New York, 1978)] 
//  
  O2cost = 0.000104 * math.Pow(v,2) + (0.182258 * v ) - 4.60
  return O2cost
}

func calcIntensity (t float64) (float64) {
// Intensity (aka "Drop Dead Formula")
// I = 0.2989558e-0.1932605t + 0.1894393e-0.012778t + 0.8
// The "t" is the amount of time in minutes that a human can run at the calculated 
// intensity expressed as a percentage of a person's maximum oxygen uptake capacity.
// Ref: Daniels, Gilbert  
  Intensity := 0.2989558 * math.Exp(-0.1932605*t) + 0.1894393 * math.Exp(-0.012778*t) + 0.8
  return Intensity
}

func calcVO2max (v, t float64) (float64) {
//
// Calculatate V02max using the Daniel's and Gilbert formula.
// The velocity is expressed in meters per minute. 
// The time is expressed in minutes. 
//
 
  O2Cost := calcO2cost(v)
  Intensity := calcIntensity(t)
  VO2max := O2Cost/Intensity
  return VO2max
}

func Bisect(fn func(float64, float64)float64, a float64, b float64, tol float64,
  maxIter int, raceLengthMeters float64) (c float64, err error) {
//
// Find the root of the function (in this case, calcVO2max) by bisecting
// https://en.wikipedia.org/wiki/Bisection_method 
// The function argument (fn) signature is specific to calcVO2max.
//
  n := 1
  for n <= maxIter {
    c := (a + b) / 2.0   // new midpoint
//    fmt.Printf("%6.2f|%6.2f|%6.2f|%6.2f|%6.2f|%6.2f\n", a,b,c, fn(c, raceLengthMeters), fn(a, raceLengthMeters), fn(10.0, 5000.0))
    if fn(c, raceLengthMeters) == 0 || (b - a)/2.0 < tol {
      return c, nil
    } else {
      n += 1
      if fn(c, raceLengthMeters) * fn(a, raceLengthMeters) > 0.0 {
	a = c
      } else {
	b = c
      }
    }
  }
  return math.NaN(), errors.New("Failed to converge.")
}

func Daniels (providedVO2Max float64, runLengthMeters float64, tStart float64, tEnd float64, 
	      raceLengthMeters float64) (tOut float64, VO2max float64, err error) {
//
// Calculate a predicted race time using the Daniel's Gilbert VO2max criteria.
// Inputs are either:
// a. A measured VO2max -or-
// b. a run length in meters and the start and end timestamps expressed 
//    as Unix-style timestamp in secs since a given reference time/date.  
//		
// Outputs are:
// tOut represents the number of minutes predicted for raceLengthMeters
// VO2max is expressed in milliliters of oxygen per kilogram of the runner's 
// weight per minute (ml/kg/min).
// err will be set if the solver failed to converge.

  if providedVO2Max == 0.0 {		
    if tStart > tEnd {
      return math.NaN(), math.NaN(), errors.New("start time after end time")
    }
    
    // Calculate the runner's VO2max based on a current run/race.
    tRun := (tEnd - tStart)/60.0 // run elapsed in mins
    vRun := runLengthMeters / tRun  // run velocity meters/min
    VO2max = calcVO2max(vRun, tRun)
  } else {
    VO2max = providedVO2Max
  }
  
  // For a race prediction, we need to solve the VO2max equation for time,
  // given a VO2max either measured or from a training run or race (above) and a 
  // distance for the race.
  // We'll use a simple bisection root solver method.
  // Implementation: What sorcery is this?  
  // Go allows function to be passed as arguments.
  // Ref: https://tour.golang.org/moretypes/24
  fn := func(raceTimeinMinutes, raceLengthMeters float64) float64 {
	    v := raceLengthMeters/raceTimeinMinutes
	    return calcVO2max(v, raceTimeinMinutes) - VO2max
	}
  a := 1.0      // Initial low guess for solution e.g. 1 minute 400 m
  b := 300.0    // Initial high guess for solution e.g. 5 hour marathon in minutes.
  tol := 0.1    // Solution tolerance 0.1 min = 6 secs = margin of error.
  maxIter := 100  //Fail if not converged after maxIter loops.
  root, err := Bisect(fn, a, b, tol, maxIter, raceLengthMeters)
  tOut = root
  return tOut, VO2max, err
  
}

func PredictRaces (providedVO2Max float64, runLengthMeters float64, tStart float64, 
		   tEnd float64) (PredictedTimes map[string]float64, VO2max float64, err error) {
//
// Make predictions.  Tell the future.  :)
//
// Predicted race times using the Daniel's Gilbert VO2max criteria.
// Inputs are either:
// a. A measured VO2max -or-
// b. a run length in meters and the start and end timestamps expressed 
//    as Unix-style timestamp in secs since a given reference time/date.  
//		
// Outputs are:
// tOut represents an array of the number of minutes predicted for common race lengths
// VO2max is expressed in milliliters of oxygen per kilogram of the runner's 
// weight per minute (ml/kg/min).
// err will be set if the solver failed to converge.

//  racelength := [6]float64 {400.0, 800.0, 5000.0, 10000.0, 21097.0, 42195.0}

  PredictedTimes = make(map[string]float64)
  PredictedTimes["400"],VO2max,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 400.0)
  PredictedTimes["800"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 800.0)
  PredictedTimes["5k"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 5000.0)
  PredictedTimes["10k"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 10000.0)
  PredictedTimes["10mi"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 16093.4)
  PredictedTimes["HM"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 21097.4)
  PredictedTimes["25k"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 25000.0)
  PredictedTimes["30k"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 30000.0)
  PredictedTimes["Mara"],_,_ = Daniels(providedVO2Max, runLengthMeters, tStart, tEnd, 42195.0)
  return PredictedTimes, VO2max, err
  }
		