#!/bin/sh
# Run Fitplot from this shell script and not Fitplot directly!
# This script makes sure the app looks in the right place for 
# associated html files.
#
# Craig Prevallet

# Change to the directory the program is run from.

cd "${0%/*}"
export PATH=$PATH:$(pwd)
./fitplot "$@"

