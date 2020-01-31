package main

import (
    "fmt"
    "image/png"
    "os"
    "github.com/gotk3/gotk3/gtk"
    "log"
)

var filename string = "test.dat"

func main() {
    // Initialize GTK without parsing any command line arguments.
    gtk.Init(nil)

    // Create a new toplevel window, set its title, and connect it to the
    // "destroy" signal to exit the GTK main loop when it is destroyed.
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle("Fitplot2")
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })

    // create a new label widget to show in the window.
    trndBtn, err := gtk.ButtonNewWithLabel("start trend")
    if err != nil {
        log.Fatal("unable to create button.", err)
    }

    // create a new label widget to show in the window.
    mapBtn, err := gtk.ButtonNewWithLabel("start map")
    if err != nil {
        log.Fatal("unable to create button.", err)
    }


    trndBtn.Connect("clicked", startTrend)
    mapBtn.Connect("clicked", startMap)

    // create a layout grid
    grid, err := gtk.GridNew()
    if err != nil {
        log.Fatal("unable to grid.", err)
    }
    grid.Attach(trndBtn, 0, 0, 100, 100)
    grid.Attach(mapBtn, 0, 101, 100, 100)


    // Add the button to the window.
    win.Add(grid)

    // Set the default window size.
    win.SetDefaultSize(50, 50)

    // Recursively show all widgets contained in this window.
    win.ShowAll()

    // Begin executing the GTK main loop.  This blocks until
    // gtk.MainQuit() is run.
    gtk.Main()
}

func startTrend() {
    InfoMessage("\nPress t to toggle scan on/off\nPress s to close window")
    go genTrendPlot(createPlotter(false,false), filename )
}
func startMap() {
        // Create and display the map.
        img := mapimg(filename)
        f, _ := os.Create("image.png")
        png.Encode(f, img)
        go displayMapPlot(createPlotter(false,false), "image.png")
}

func InfoMessage(format string, args ...interface{}) error {
	dummy, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return err
	}
	dialog := gtk.MessageDialogNew(dummy, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, fmt.Sprint("INFO: ", format), args...)
	dialog.Run()
	dialog.Destroy()
	return nil
}
