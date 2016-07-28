#!/bin/bash
#
# Package Fitplot for various operating systems.
# Usage package {help|windows|linux|osx|all}
# 
#

#
# Useful packaging tools assumed available:
# 
# iconify2 gimp plugin http://registry.gimp.org/node/27989 png->.ico
# github.com/akavel/rsrc/ - embed windows icons - go get github.com/akavel/rsrc
# nsis - windows installer utility - apt-get install nsis
# hfsplus + dependencies - osx file system - apt-get install hfsplus
 
# Configure these to match your system
SOURCE_DIR=/home/penguin/work/src/github.com/cprevallet/fitplot
BUILD_DIR=/home/penguin/work/builds
MOUNT_DIR=/mnt/fitplot

function build() {
	local os="$1"
	local arch="$2"
	cd $SOURCE_DIR/static
	rst2html --stylesheet=help.css  help.rst > help.html
	rst2odt help.rst > help.odt
	rst2man help.rst > help.man
	# The following must be done manually:
	# In gimp open fitplot.xcf and use iconify plugin to create windows object file.
	# rsrc -ico fitplot_color.ico -o fitplot.syso -arch amd64
	# export fitplot.xcf to fitplot_color.png
	# Create pdf file from help.rst
	# Open help.odt in LibreOffice, right click on table of contents and update - generates page numbers
	# export help.odt to pdf
	#
	
	cd $SOURCE_DIR
    echo -e '\nbuilding:'$os, $arch
	go clean
	if [ "$os" == "windows" ] ; then
		env GOOS=$1 GOARCH=$2 go build -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -H windowsgui" -v
	else
		env GOOS=$1 GOARCH=$2 go build -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" -v
	fi
}

function package_windows {
    echo -e 'Packaging for Windows'
	cd $BUILD_DIR
	# make a backup of previous package
	sudo cp -r windows_dist/ windows_dist_old
	sudo rm -rf windows_dist
	sudo mkdir windows_dist
	sudo mkdir windows_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot.exe ./windows_dist/fitplot
	sudo cp $SOURCE_DIR/LICENSE.txt ./windows_dist/fitplot
	sudo cp -r $SOURCE_DIR/static/ ./windows_dist/fitplot
	sudo cp -r $SOURCE_DIR/tmpl/ ./windows_dist/fitplot
	sudo cp -r $SOURCE_DIR/samples/ ./windows_dist/fitplot
	cd windows_dist
	sudo cp ../Fitplot\ Windows\ x64\ Setup.nsi .
	sudo makensis Fitplot\ Windows\ x64\ Setup.nsi
	sudo sh -c "md5sum Fitplot\ Windows\ x64\ Setup.exe > md_windows.txt"
}

function package_linux {
    echo -e 'Packaging for Linux'
	cd $BUILD_DIR
	# make a backup of previous package
	sudo cp -r linux_dist/ linux_dist_old
	sudo rm -rf linux_dist
	sudo mkdir linux_dist
	sudo mkdir linux_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot ./linux_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot.desktop ./linux_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot.sh ./linux_dist/fitplot
	sudo cp -r $SOURCE_DIR/static/ ./linux_dist/fitplot
	sudo cp -r $SOURCE_DIR/tmpl/ ./linux_dist/fitplot
	sudo cp -r $SOURCE_DIR/samples/ ./linux_dist/fitplot
	sudo cp -r $SOURCE_DIR/icons/ ./linux_dist/fitplot
	cd linux_dist
	sudo chown root:root fitplot -R
	sudo sh -c "tar -cvzf fitplot_linux64bit.tgz fitplot/"
	sudo sh -c "md5sum fitplot_linux64bit.tgz > md_linux.txt"
}

function package_osx {
	echo -e 'Packaging for OSX'
	cd $BUILD_DIR
	# make a backup of previous package
	sudo cp -r osx_dist/ osx_dist_old
	sudo rm -rf osx_dist
	sudo mkdir osx_dist
	sudo mkdir osx_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot ./osx_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot.desktop ./osx_dist/fitplot
	sudo cp $SOURCE_DIR/fitplot.sh ./osx_dist/fitplot
	sudo cp -r $SOURCE_DIR/static/ ./osx_dist/fitplot
	sudo cp -r $SOURCE_DIR/tmpl/ ./osx_dist/fitplot
	sudo cp -r $SOURCE_DIR/samples/ ./osx_dist/fitplot
	sudo cp -r $SOURCE_DIR/icons/ ./osx_dist/fitplot
	cd osx_dist
	sudo chown root:root fitplot -R
	sudo dd if=/dev/zero of=fitplot_osx64bit.dmg bs=1M count=20
	sudo mkfs.hfsplus -v Fitplot fitplot_osx64bit.dmg
	# Does mount directory exist?
#	if [! -d  "$MOUNT_DIR"]; then
#		sudo mkdir /mnt/fitplot
#	fi
	sudo mount -o loop fitplot_osx64bit.dmg $MOUNT_DIR
	cd fitplot
	sudo mv fitplot.sh fitplot.command
	sudo cp -r .  /mnt/fitplot
	sudo umount $MOUNT_DIR
	cd ..
	sudo sh -c "md5sum fitplot_osx64bit.dmg > md_osx.txt"
}

function build_and_package {
    echo -e "$1"
    if [ "$1" = 'windows' ]; then
		build "windows" "amd64"
        package_windows
    fi
    if [ "$1" = 'linux' ]; then
		build "linux" "amd64"
        package_linux
    fi
    if [ "$1" = 'osx' ]; then
		build "darwin" "amd64"
        package_osx
    fi
    if [ "$1" = 'all' ]; then
		build "windows" "amd64"
        package_windows
		build "linux" "amd64"
        package_linux
		build "darwin" "amd64"
        package_osx
    fi
}



# Main entry point
# First arg is help or no args passed.
if [ "$1" = 'help' ] || [ $# -eq 0 ]; then
    echo -e '\nUsage: package {help|windows|linux|osx|all}'
else
    date "+%nPackage started on date: %m-%d-%Y at %H:%M:%S" 
    until [ -z "$1" ]  # Until all parameters used up . . .
    do
      build_and_package $1
      shift
    done
    date "+%nPackage ended on date: %m-%d-%Y at %H:%M:%S"
fi

