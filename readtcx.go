package main

import (
	"fmt"
        "github.com/csp/fitplot/tcx"
	"github.com/jezard/fit"
)

var tcxFname string = "./tcx/sample.tcx"

func printTCX() {
  db, _ := tcx.ReadTCXFile(tcxFname)
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
	
  //First convert to a structure we already know how to handle!
  //Then we can re-use the existing routines in graphdata.go
  myStruct := cvtTofitStruct(db)
  for _, record := range myStruct.Records {
      fmt.Println(record)
  }
  
  //TODO call getDvsA, getDvsC, getDvsP, getlatlong
}

func cvtTofitStruct(db *tcx.TCXDB) (fitStruct fit.FitFile) {
  for i, _ := range db.Acts.Act {
	  for j, _ := range db.Acts.Act[i].Laps {
		  for k, _ := range db.Acts.Act[i].Laps[j].Trk.Pt {                   
		    var newRec fit.Record
		    newRec.Position_lat = db.Acts.Act[i].Laps[j].Trk.Pt[k].Lat
		    newRec.Position_long = db.Acts.Act[i].Laps[j].Trk.Pt[k].Long
		    newRec.Altitude = db.Acts.Act[i].Laps[j].Trk.Pt[k].Alt
		    newRec.Distance = db.Acts.Act[i].Laps[j].Trk.Pt[k].Dist
		    newRec.Speed = db.Acts.Act[i].Laps[j].Trk.Pt[k].Speed
		    newRec.Cadence = uint8(db.Acts.Act[i].Laps[j].Trk.Pt[k].Cad)
		    fitStruct.Records = append(fitStruct.Records, newRec)
		  }
	  }
  }
  return fitStruct
}