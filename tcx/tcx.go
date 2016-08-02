package tcx

import (
	"encoding/xml"
	"io/ioutil"
)

// ReadTpts converts the GPS track points from XML.
func ReadTpts(path string) (track *Track, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	track = new(Track)
	err = xml.Unmarshal(filebytes, track)
	return
}

// ReadLap converts the lap values from XML.
func ReadLap(path string) (lap *Lap, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lap = new(Lap)
	err = xml.Unmarshal(filebytes, lap)
	return
}

// ReadActivity converts an individual activity value from XML.
func ReadActivity(path string) (act *Activity, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	act = new(Activity)
	err = xml.Unmarshal(filebytes, act)
	return
}

// ReadActivities converts an array of activities values from XML.
func ReadActivities(path string) (acts *Activities, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	acts = new(Activities)
	err = xml.Unmarshal(filebytes, acts)
	return
}

// ReadTCXFile converts an entire TCX file from XML.
func ReadTCXFile(path string) (db *TCXDB, err error) {
	filebytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	db = new(TCXDB)
	err = xml.Unmarshal(filebytes, db)
	return
}
