package main

import (
        "github.com/cprevallet/fitplot/tcx"
	"github.com/jezard/fit"
	"math"
	"time"
)

// Convert the TCXDB structure created from the XML to fit.Record structure
func cvtToFitRecs(db *tcx.TCXDB) (runRecs []fit.Record) {
  

  var hasspeed bool
  var hasdist bool
 
  // Determine if TCX file supplied cumulative distance and speed.
    for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                   
		    if db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed != 0.0 {
		      hasspeed = true
		    }
		    if db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist != 0.0 {
		      hasdist = true
		    }
		  }  
	  }
    }
 
  // Calculate cumulative distance in meters using latitude and longitude if not supplied in XML.
  if (!hasdist) {
     var lat0, long0, lat1, long1, totalDist, newDist float64
     for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                    
		    if (i==0 && j==0 && k==0) {
			lat0 = db.Acts.Act[0].Laps[0].Trk.Pt[0].Lat
			long0 = db.Acts.Act[0].Laps[0].Trk.Pt[0].Long
		      } else {
			lat1 = db.Acts.Act[i].Laps[j].Trk.Pt[k].Lat
			long1 = db.Acts.Act[i].Laps[j].Trk.Pt[k].Long
			newDist = Distance(lat0, long0, lat1, long1)
			totalDist = totalDist + newDist
			lat0 = lat1
			long0 = long1
			db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist = totalDist
		      }
		  }
	  }
    }
  }
  
  // Calculate segment speed in meters/sec using distance and timestamp if not supplied in XML.
  if (!hasspeed) {
      var lasttime, thistime time.Time
      var deltaT time.Duration
      var dist, lastdist float64  //meters
         for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                    
		    if (i==0 && j==0 && k==0) {
		      lasttime = db.Acts.Act[i].Laps[j].Trk.Pt[k].Time
		      lastdist = 0.0
		      db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed = 0.0
		  } else {
		      thistime = db.Acts.Act[i].Laps[j].Trk.Pt[k].Time
		      dist = db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist
		      deltaT = thistime.Sub(lasttime)
		      if deltaT.Seconds() > 0.01 {
			db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed = (dist-lastdist)/deltaT.Seconds()
		      } else {
			db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed = 0.0
		      }
		      lasttime = thistime
		      lastdist = dist
		    }
		    }
	  }
    }
  }
  
  for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                   
		    
		    var newRec fit.Record
		    newRec.Timestamp = db.Acts.Act[i].Laps[j].Trk.Pt[k].Time.Unix()
		    newRec.Position_lat = db.Acts.Act[i].Laps[j].Trk.Pt[k].Lat
		    newRec.Position_long = db.Acts.Act[i].Laps[j].Trk.Pt[k].Long
		    newRec.Altitude = db.Acts.Act[i].Laps[j].Trk.Pt[k].Alt
		    newRec.Distance = db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist
		    newRec.Speed = db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed
		    newRec.Cadence = uint8(db.Acts.Act[i].Laps[j].Trk.Pt[k].Cad)
		    if newRec.Position_lat != 0.0 && newRec.Position_long != 0.0 {
		      runRecs = append(runRecs, newRec)
		    }
		  }
	  }
  }
  return runRecs
}

// Convert the TCXDB structure created from the XML to fit.Laps structure
func cvtToFitLaps(db *tcx.TCXDB) (runLaps []fit.Lap) {

    for i, _ := range db.Acts.Act {
	  for _,lap := range db.Acts.Act[i].Laps {
		    var newLap fit.Lap
		    //TODO This doesn't seem to work.  Maybe time.RFC3339 isn't right...
		    //t, err :=  time.Parse(lap.Start, time.RFC3339)
		    //if err != nil {newLap.Timestamp = t.Unix()} else {break}
		    newLap.Total_elapsed_time = lap.TotalTime
		    newLap.Total_distance = lap.Dist
		    newLap.Total_calories = uint16(lap.Calories)
		    runLaps = append(runLaps, newLap)
		    }
		  }
     return runLaps
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	var la1, lo1, la2, lo2, r, dlat, dlon float64
	la1 = lat1 * math.Pi / 180.0
	lo1 = lon1 * math.Pi / 180.0
	la2 = lat2 * math.Pi / 180.0
	lo2 = lon2 * math.Pi / 180.0

	r = 6378100.0 // Earth radius in METERS

	// find the differences between the coordinates
	dlat = la2 - la1
	dlon = lo2 - lo1
		
	// here's the heavy lifting
	a  := math.Pow(math.Sin(dlat/2.0),2) + math.Cos(la1) * math.Cos(la2) * math.Pow(math.Sin(dlon/2.0),2)
	c  := 2.0 * math.Atan2(math.Sqrt(a),math.Sqrt(1.0-a)) // great circle distance in radians	
	
	return c * r
}


// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}