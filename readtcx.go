package main

import (
        "github.com/cprevallet/fitplot/tcx"
	"github.com/jezard/fit"
	"math"
)

func readTCX(tcxFname string) (runRecs []fit.Record) {
  
  db, _ := tcx.ReadTCXFile(tcxFname)
  /*
	  for i, _ := range db.Acts.Act {
		  for j, _ := range db.Acts.Act[i].Laps {
			  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                   
			  fmt.Println("Latitude", db.Acts.Act[i].Laps[j].Trk.Pt[k].Lat)       
			  fmt.Println("Longitude", db.Acts.Act[i].Laps[j].Trk.Pt[k].Long)     
			  fmt.Println("Altitude (m)", db.Acts.Act[i].Laps[j].Trk.Pt[k].Alt)   
			  fmt.Println("Distance (m)", db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist)  
			  fmt.Println("Speed (m/s)", db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed)  
			  fmt.Println("Cadence (bpm)", db.Acts.Act[i].Laps[j].Trk.Pt[k].Cad)  
			  }
		  }
	}
*/

  //Convert to a structure we already know how to handle!
  //Then we can re-use the existing routines in graphdata.go.
  return cvtToFitRecs(db)
  
}

// Convert the TCXDB structure created from the XML to fit.Record structure
func cvtToFitRecs(db *tcx.TCXDB) (runRecs []fit.Record) {
  var startLat float64
  var startLong float64
  for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                   
		    var newRec fit.Record
		    newRec.Position_lat = db.Acts.Act[i].Laps[j].Trk.Pt[k].Lat
		    newRec.Position_long = db.Acts.Act[i].Laps[j].Trk.Pt[k].Long
		    newRec.Altitude = db.Acts.Act[i].Laps[j].Trk.Pt[k].Alt
		    // XML has no totalized run distance! Must calculate with lat/long.
		    if (i==0 && j==0 && k==0) {
		      newRec.Distance = float64(0.0)
		      startLat = newRec.Position_lat
		      startLong = newRec.Position_long
		    } else {
		      newRec.Distance = Distance(newRec.Position_lat, newRec.Position_long, startLat, startLong)
		    }
		    //newRec.Distance = db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist
		    newRec.Speed = db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed
		    // TODO the following typecast is generating zeros?
		    newRec.Cadence = uint8(db.Acts.Act[i].Laps[j].Trk.Pt[k].Cad)
		    //fmt.Println(newRec.Cadence, db.Acts.Act[i].Laps[j].Trk.Pt[k].Cad)
		    runRecs = append(runRecs, newRec)
		  }
	  }
  }
  return runRecs
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