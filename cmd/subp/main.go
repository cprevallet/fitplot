package main

import "fmt"

//import "io/ioutil"
import "github.com/sbinet/go-gnuplot"

func main() {
	fname := ""
	persist := true 
	debug := true

	p, err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p2, err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p2.Close()

	p.CheckedCmd("set terminal wxt size 1024,800")
        p.CheckedCmd("set multiplot layout 3,1")
        p.CheckedCmd("set border linewidth 1.5")
	p.CheckedCmd("set datafile separator ','")
        p.CheckedCmd("set style line 1 linecolor rgb '#0060ad' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5")
        p.CheckedCmd("set style line 2 linecolor rgb '#dd181f' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5")
        p.CheckedCmd("set style line 3 linecolor rgb '#5416b4' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5")
        p.CheckedCmd("set grid")
	p.CheckedCmd("unset key")
        p.CheckedCmd("set title 'Pace'")
        p.CheckedCmd("set ylabel 'Pace, min/km'")
        p.CheckedCmd("set ydata time")
        p.CheckedCmd("unset xlabel")
        p.CheckedCmd("set timefmt %s", "'%M:%S'")
        p.CheckedCmd("set yrange [*:*] reverse")
        p.CheckedCmd("plot 'test.dat' using 1:3  with linespoints linestyle 1")
        p.CheckedCmd("unset key")
        p.CheckedCmd("unset ydata")
        p.CheckedCmd("unset timefmt")
        p.CheckedCmd("unset yrange")
        p.CheckedCmd("set title 'Cadence'")
        p.CheckedCmd("set ylabel 'Strides/min'")
        p.CheckedCmd("unset xlabel")
        p.CheckedCmd("plot 'test.dat' using 1:7  with linespoints linestyle 2")
        p.CheckedCmd("set title 'Altitude'")
        p.CheckedCmd("set ylabel 'Altitude, m'")
        p.CheckedCmd("set xlabel 'Distance, m'")
        p.CheckedCmd("plot 'test.dat' using 1:6  with linespoints linestyle 3")
	p2.CheckedCmd("set terminal wxt size 1024,800")
        p2.CheckedCmd("unset xlabel")
        p2.CheckedCmd("unset ylabel")
        p2.CheckedCmd("unset xtics")
        p2.CheckedCmd("unset ytics")
        p2.CheckedCmd("set size ratio -1")
        p2.CheckedCmd("plot 'image.png' binary filetype=png with rgbimage")
        p2.CheckedCmd("pause 1")

	//p.CheckedCmd("q")
	//p.proc.Wait(0)

	return
}

