#!/bin/sh
# Run Fitplot from this shell script and not Fitplot directly!
# This script makes sure the app looks in the right place for 
# associated html files.
#
# Craig Prevallet

# Change to the directory the program is run from.

#cd "${0%/*}"
#export PATH=$PATH:$(pwd)
#./fitplot "$@"

echo "Fitplot has extracted itself."
echo "Setting up desktop and icons files."
cp ./nw.package/fitplot.desktop /usr/share/applications 
cp ./nw.package/icons/fitplot_color.png /usr/share/pixmaps/fitplot.png
echo "All done."
