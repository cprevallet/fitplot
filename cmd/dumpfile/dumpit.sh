# the following line is needed to get the scrolling to work under Wayland
export GDK_BACKEND=x11
../../dumpfile INPUTDATA.FIT > test.dat ; gnuplot plotdata.gnu -p
rm test.dat
