package main

import (
  "bufio"
  "encoding/csv"
  "image"
  "image/color"
  //"image/png"
  "io"
  "github.com/flopp/go-staticmaps"
  "github.com/golang/geo/s2"
  "os"
  "strconv"
)


func readInputs() ([]s2.LatLng) {
    csvFile, _ := os.Open("test.dat")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    var position []s2.LatLng
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            panic(error)
        }
        lat, err := strconv.ParseFloat(line[3], 64)
        if err != nil {panic(err)}
        lng, err := strconv.ParseFloat(line[4], 64)
        if err != nil {panic(err)}
        position = append(position, s2.LatLngFromDegrees(lat, lng))
    }
    return position
}


func mapimg() (image.Image) {
  position := readInputs()
  ctx := sm.NewContext()
  ctx.SetSize(800, 600)
  ctx.AddMarker(sm.NewMarker(position[0], color.RGBA{0xff, 0, 0, 0xff}, 16.0))
  ctx.AddPath(sm.NewPath(position, color.RGBA{0xff, 0, 0, 0xff}, 2.0))
  ctx.AddMarker(sm.NewMarker(position[len(position)-1], color.RGBA{0xff, 0, 0, 0xff}, 16.0))

  img, err := ctx.Render()
  if err != nil {
    panic(err)
  }
  return img
}
