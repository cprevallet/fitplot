# the following line is needed to get the scrolling to work under Wayland
export GDK_BACKEND=x11
cd ../dumpfile/
../../dumpfile > test.dat ; gnuplot ../dumpfile/plotdata.gnu -p & 
cd ../mapfile/
cat ../dumpfile/test.dat | ../../mapfile; eog image.png &
