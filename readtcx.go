package main

import (
	"fmt"
        "github.com/csp/fitplot/tcx"
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
}
