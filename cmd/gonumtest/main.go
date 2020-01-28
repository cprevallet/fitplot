// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
        "bufio"
        "encoding/csv"
        "io"
        "os"
	"image/color"
        //"fmt"
	"log"
        "strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"

        "github.com/cprevallet/fitplot/strutil"
)

type customTicks struct{}
// Ticks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
func (customTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		tks[i].Label = customizeLabels(t.Label)
	}
	return tks
}

// AddCommas adds commas after every 3 characters from right to left.
// NOTE: This function is a quick hack, it doesn't work with decimal
// points, and may have a bunch of other problems.
func customizeLabels(s string) string {
        f, _ := strconv.ParseFloat(s, 64)
	return strutil.DecimalTimetoMinSec(f)
}

var minPace = 0.56
var paceToMetric = 16.666667 // sec/meter -> min/km

func makePace( speed float64) ( pace float64) {
	// Speed -> pace
	if speed > 1.8 { //m/s
		pace = 1.0/speed
	} else {
		pace = minPace // s/m = 15 min/mi
	}
        return
        }

func readInputs() plotter.XYs  {
    csvFile, _ := os.Open("test.dat")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    var pts plotter.XYs
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            panic(error)
        }
        x, err := strconv.ParseFloat(line[0], 64)
        if err != nil {panic(err)}
        y, err := strconv.ParseFloat(line[1], 64)
        if err != nil {panic(err)}
        pts = append(pts, plotter.XY {X:x, Y:makePace(y)*paceToMetric})  
    }
    return pts
}

// ExampleScatter draws some scatter points, a line,
// and a line with points.
func main() {
	linePointsData := readInputs()

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Pace Graph"
	p.X.Label.Text = "Distance, m"
	p.Y.Label.Text = "Pace, min/km"

        // Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks.
	p.Y.Tick.Marker = customTicks{}

        p.Y.Scale = plot.InvertedScale{Normalizer: plot.LinearScale{}}
	p.Add(plotter.NewGrid())

	if err != nil {
		log.Panic(err)
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(linePointsData)
	if err != nil {
		log.Panic(err)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.CircleGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(lpLine, lpPoints)

	err = p.Save(800, 600, "scatter.png")
	if err != nil {
		log.Panic(err)
	}
}
