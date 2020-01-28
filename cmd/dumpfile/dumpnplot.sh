# the following line is needed to get the scrolling to work under Wayland
export GDK_BACKEND=x11
../../dumpfile $1 > test.dat ; ../../gplot test.dat
rm test.dat
