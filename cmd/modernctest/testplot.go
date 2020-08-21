package main

import (
        "fmt"
        "modernc.org/plot"
       )

const src = `
set terminal wxt 1 size 1024,800 persist
bind 'x' 'unset terminal; exit gnuplot'
set border linewidth 1.5
set datafile separator ','
set style line 1 linecolor rgb '#0060ad' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set style line 2 linecolor rgb '#dd181f' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set style line 3 linecolor rgb '#5416b4' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set grid xtics ytics
set title 'Pace'
set ydata time
unset xlabel
set timefmt '%M:%S'
set yrange [*:*] reverse
unset y2range
set autoscale y2
set y2tics
set tics out
set title 'Pace Graph'
set ylabel 'Pace, min/km'
set y2label 'Altitude, m'
set xlabel 'Distance, m'
plot 'test.dat' using 1:3  with linespoints linestyle 1 axes x1y1 title 'Pace','test.dat' using 1:6  with linespoints linestyle 3 axes x1y2 title 'Altitude'
mycond = 1
while (mycond == 1) {}
`

const src2 = `
set terminal wxt 1 size 1024,800 persist
bind 'x' 'unset terminal; exit gnuplot'
set border linewidth 1.5
set datafile separator ','
set style line 1 linecolor rgb '#0060ad' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set style line 2 linecolor rgb '#dd181f' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set style line 3 linecolor rgb '#5416b4' linetype 1 linewidth 0.5 pointtype 7 pointsize 0.5
set grid xtics ytics
set title 'Pace'
set ydata time
unset xlabel
set timefmt '%M:%S'
set yrange [*:*] reverse
unset y2range
set autoscale y2
set y2tics
set tics out
set title 'Cadence Graph'
set ylabel 'Pace, min/km'
set y2label 'Cadence, strides/m'
set xlabel 'Distance, m'
plot 'test.dat' using 1:3  with linespoints linestyle 1 axes x1y1 title 'Pace','test.dat' using 1:7  with linespoints linestyle 2 axes x1y2 title 'Cadence'
mycond = 1
while (mycond == 1) {}
`

func main () {
        fmt.Printf("%s\n", src)
        // Create plots using the gnuplot API
        out, err := plot.Script([]byte(src))
        if err != nil {
            panic(err)
        }
        fmt.Printf("%s\n", out)
        fmt.Printf("%s\n", src2)
        // Create plots using the gnuplot API
        out, err = plot.Script([]byte(src2))
        if err != nil {
            panic(err)
        }
        fmt.Printf("%s\n", out)

}
