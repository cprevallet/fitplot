package tcx

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Trackpoint represents an individual GPS track and associated measured 
// values at a particular time.
type Trackpoint struct {
	Time  time.Time
	Lat   float64 `xml:"Position>LatitudeDegrees"`
	Long  float64 `xml:"Position>LongitudeDegrees"`
	Alt   float64 `xml:"AltitudeMeters,omitempty"`
	Dist  float64 `xml:"DistanceMeters,omitempty"`
	HR    float64 `xml:"HeartRateBpm>Value,omitempty"`
	Cad   float64 `xml:"Cadence,omitempty"`
	Speed float64 `xml:"Extensions>TPX>Speed,omitempty"`
	Power float64 `xml:"Extensions>TPX>Watts,omitempty"`
}

// Lap represents summary values for GPS and activity related information.
type Lap struct {
	Start         string  `xml:"StartTime,attr"`
	TotalTime     float64 `xml:"TotalTimeSeconds,omitempty"`
	Dist          float64 `xml:"DistanceMeters,omitempty"`
	Calories      float64 `xml:",omitempty"`
	MaxSpeed      float64 `xml:"MaximumSpeed,omitempty"`
	AvgHr         float64 `xml:"AverageHeartRateBpm>Value,omitempty"`
	MaxHr         float64 `xml:"MaximumHeartRateBpm>Value,omitempty"`
	Intensity     string  `xml:",omitempty"`
	TriggerMethod string  `xml:",omitempty"`
	Trk           *Track  `xml:"Track"`
}

// Track represents an array of Trackpoint.
type Track struct {
	Pt []Trackpoint `xml:"Trackpoint"`
}

// Activity represents an individual sporting activity undertaken at a 
// specific time, the device it was captured with and the associated 
// raw summary (lap) information.
type Activity struct {
	Sport   string `xml:"Sport,attr,omitempty"`
	Id      time.Time
	Laps    []Lap  `xml:"Lap,omitempty"`
	Creator Device `xml:"Creator,omitempty"`
}

// Device respresents the hardware used to capture information.
type Device struct {
	Name      string       `xml:",omitempty"`
	UnitId    uint         `xml:",omitempty"`
	ProductID string       `xml:",omitempty"`
	Version   BuildVersion `xml:",omitempty"`
}

// TCXDB is the top level file structure containing Activities.
type TCXDB struct {
	XMLName xml.Name    `xml:"TrainingCenterDatabase"`
	Acts    *Activities `xml:"Activities"`
	Auth    Author      `xml:"Author"`
}

// Activities represents an array of Activity.
type Activities struct {
	Act []Activity `xml:"Activity"`
}

// Author represents the software used on a Device to capture information.
type Author struct {
	Name       string `xml:",omitempty"`
	Build      Build  `xml:",omitempty"`
	LangID     string `xml:",omitempty"`
	PartNumber string `xml:",omitempty"`
}

// Function to retrieve information about the software.
func (a Author) String() string {
	return fmt.Sprintf("%v: version %v.%v.%v.%v", a.Name, a.Build.Version.VersionMajor, a.Build.Version.VersionMinor, a.Build.Version.BuildMajor, a.Build.Version.BuildMinor)
}

// Build is meta-information about the Garmin software.
type Build struct {
	Version BuildVersion `xml:"Version,omitempty"`
	Type    string       `xml:",omitempty"`
	Time    string       `xml:",omitempty"`
	Builder string       `xml:",omitempty"`
}

// BuildVersion is meta-information about the Garmin software version.
type BuildVersion struct {
	VersionMajor int `xml:",omitempty"`
	VersionMinor int `xml:",omitempty"`
	BuildMajor   int `xml:",omitempty"`
	BuildMinor   int `xml:",omitempty"`
}
