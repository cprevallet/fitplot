package strutil

//
// Open a fit file and organize as datapoints suitable to plot.
//

import (
//  "fmt"
  "math"
  "strconv"
)

// Convert time in decimal minutes to mm:ss.
// Does not handle input values > one hour.
// See function decimalTimetoHourMinSec.
func DecimalTimetoMinSec(in float64) (out string) {
  in_min := int(math.Floor(in))
  in_min_str := strconv.Itoa(in_min)
  if in_min < 10 { in_min_str = "0" + in_min_str}
  in_sec := int((in - float64(in_min))* 60)
  in_sec_str := strconv.Itoa(in_sec)
  if in_sec < 10 { in_sec_str = "0" + in_sec_str}
  out = in_min_str + ":" + in_sec_str
  return out
}

// Convert decimal time in mins to hh:mm:ss
func DecimalTimetoHourMinSec(in float64) (out string) {
  totalsecs := float64(in * 60.0) //seconds
  
  hours := math.Floor(totalsecs / 3600.0)
  hours_str := strconv.Itoa(int(hours))
  if hours < 10 {hours_str = "0" + hours_str}
  totalsecs -= hours * 3600
  
  mins := math.Floor(totalsecs / 60.0)
  mins_str := strconv.Itoa(int(mins))
  if mins < 10 {mins_str = "0" + mins_str}
  totalsecs -= mins * 60
  
  secs := math.Floor(totalsecs)
  secs_str := strconv.Itoa(int(secs))
  if secs < 10 {secs_str = "0" + secs_str}
  out = hours_str + ":" + mins_str + ":" + secs_str 
  return out
}