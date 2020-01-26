// Package main provides various examples of Fyne API capabilities
package main

import (
//	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Demo")

        img := mapimg()  //image.Image - see genmap.go file

	logo := canvas.NewImageFromImage(img)
	logo.SetMinSize(fyne.NewSize(600, 600))

	w.SetContent(widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
			a.Quit()
		}),
	))
	w.ShowAndRun()
}
