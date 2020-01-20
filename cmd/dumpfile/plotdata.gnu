#!/usr/bin/gnuplot
#
# Plotting the data of file plotting_data3.dat
#
# AUTHOR: Craig Prevallet
# 
# run via: cat data.txt | gnuplot -persist plotdata.gnu

set term wxt
#set term tkcanvas
#set term x11 
#reset
set title "Test Plot"
set border linewidth 1.5
# Set first two line styles to blue (#0060ad) and red (#dd181f)
set style line 1 \
    linecolor rgb '#0060ad' \
    linetype 1 linewidth 2 \
    pointtype 7 pointsize 1.5
set style line 2 \
    linecolor rgb '#dd181f' \
    linetype 1 linewidth 2 \
    pointtype 5 pointsize 1.5

unset key
set ylabel "Speed m/s"
set xlabel "Distance, m"
# set ydata time
# set timefmt "%M:%S"
# set yrange [*:*] reverse
set datafile separator ','
plot 'test.dat' using 1:2  with linespoints
pause 10
reread
