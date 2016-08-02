//
// Package strutil provides conversion routines from decimal times to their
// equivalent strings. (e.g. 9.5 to "9:30")
//
package strutil

//
// Open a fit file and organize as datapoints suitable to plot.
//

import (
	//"fmt"
	"math"
	"strconv"
)

// DecimalTimetoMinSec converts a time in floating point decimal minutes to
// a string of mm:ss. (e.g. 9.5 to "9:30")
// Does not handle input values > one hour. See function decimalTimetoHourMinSec.
func DecimalTimetoMinSec(in float64) (out string) {
	inMin := int(math.Floor(in))
	inMinStr := strconv.Itoa(inMin)
	if inMin < 10 {
		inMinStr = "0" + inMinStr
	}
	inSec := int((in - float64(inMin)) * 60)
	inSecStr := strconv.Itoa(inSec)
	if inSec < 10 {
		inSecStr = "0" + inSecStr
	}
	out = inMinStr + ":" + inSecStr
	return out
}

// DecimalTimetoHourMinSec converts a time in floating point decimal minutes to
// a string of hh:mm:ss (e.g. 67.5 to "01:07:30").
func DecimalTimetoHourMinSec(in float64) (out string) {
	totalSecs := float64(in * 60.0) //seconds

	hours := math.Floor(totalSecs / 3600.0)
	hoursStr := strconv.Itoa(int(hours))
	if hours < 10 {
		hoursStr = "0" + hoursStr
	}
	totalSecs -= hours * 3600

	mins := math.Floor(totalSecs / 60.0)
	minsStr := strconv.Itoa(int(mins))
	if mins < 10 {
		minsStr = "0" + minsStr
	}
	totalSecs -= mins * 60

	secs := math.Floor(totalSecs)
	secsStr := strconv.Itoa(int(secs))
	if secs < 10 {
		secsStr = "0" + secsStr
	}
	out = hoursStr + ":" + minsStr + ":" + secsStr
	return out
}
