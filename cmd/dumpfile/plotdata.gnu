#!/usr/bin/gnuplot
#
# Plotting the data of file plotting_data3.dat
#
# AUTHOR: Craig Prevallet
# 
# run via: cat data.txt | gnuplot -persist plotdata.gnu


if (! exists("setsize")) set term wxt size 1024,800;
setsize = 1
#set term wxt size 1024,800
#set term tkcanvas
#set term x11 
#reset
# ctrl-x = quit
bind "ctrl-x" "exit gnuplot"

set multiplot layout 3,1 
set border linewidth 1.5
# Set first two line styles to blue (#0060ad) and red (#dd181f)
set style line 1 \
    linecolor rgb '#0060ad' \
    linetype 1 linewidth 0.5 \
    pointtype 7 pointsize 0.5
set style line 2 \
    linecolor rgb '#dd181f' \
    linetype 1 linewidth 0.5 \
    pointtype 7 pointsize 0.5
set style line 3 \
    linecolor rgb '#5416b4' \
    linetype 1 linewidth 0.5 \
    pointtype 7 pointsize 0.5
set grid
unset key
set title "Pace"
set ylabel "Pace min/km"
#set xlabel "Distance, m"
unset xlabel
set ydata time
set timefmt "%M:%S"
set yrange [*:*] reverse
set datafile separator ','
plot 'test.dat' using 1:3  with linespoints linestyle 1

unset key
unset ydata
unset timefmt
unset yrange

set title "Cadence"
set ylabel "Strides/min"
unset xlabel
#set xlabel "Distance, m"
set datafile separator ','
plot 'test.dat' using 1:7  with linespoints linestyle 2

set title "Altitude"
set ylabel "Altitude m"
set xlabel "Distance, m"
set datafile separator ','
plot 'test.dat' using 1:6  with linespoints linestyle 3 
#unset multiplot
pause 2
replot
reread
