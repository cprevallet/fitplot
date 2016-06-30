package tcx

import (
	"encoding/xml"
	"io/ioutil"
)

// Read the track points.
func ReadTpts(path string) (track *Track, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	track = new(Track)
	err = xml.Unmarshal(filebytes, track)
	return
}

// Read the lap.
func ReadLap(path string) (lap *Lap, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lap = new(Lap)
	err = xml.Unmarshal(filebytes, lap)
	return
}

// Read the activity.
func ReadActivity(path string) (act *Activity, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	act = new(Activity)
	err = xml.Unmarshal(filebytes, act)
	return
}

// Read the the activities.
func ReadActivities(path string) (acts *Activities, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	acts = new(Activities)
	err = xml.Unmarshal(filebytes, acts)
	return
}

// Read the TCX file.
func ReadTCXFile(path string) (db *TCXDB, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	db = new(TCXDB)
	err = xml.Unmarshal(filebytes, db)
	return
}
