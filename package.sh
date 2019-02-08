#!/bin/bash
#
# Package Fitplot for various operating systems.
# Usage package {help|windows|linux|osx|all}
# 
#

#
# Useful packaging tools assumed available:
#
# asciidoctor - apt-get install asciidoctor
# iconify2 gimp plugin http://registry.gimp.org/node/27989 png->.ico
# github.com/akavel/rsrc/ - embed windows icons - go get github.com/akavel/rsrc
# nsis - windows installer utility - apt-get install nsis
# hfsplus + dependencies - osx file system - apt-get install hfsplus
# cross-compilers for C language programs (e.g. sqlite3 dependency is cgo )
# x86_64-w64-mingw32-gcc and o64-clang
# see: https://www.limitlessfx.com/cross-compile-golang-app-for-windows-from-linux.html
#      https://github.com/tpoechtrager/osxcross
 
# Configure these to match your system
BIN_DIR=/home/craig/go/bin
SOURCE_DIR=/home/craig/go/src/github.com/cprevallet/fitplot
PKG_DIR=/home/craig/go/packaging
MOUNT_DIR=/mnt/fitplot
NWJS_DIR=/home/craig/Downloads/

function build() {
	local os="$1"
	local arch="$2"
	cd $SOURCE_DIR/static
        asciidoctor -a toc -a icons=font -a icon-set=fa help.asc
        asciidoctor-pdf -a toc -a icons=font -a icon-set=fa help.asc
	# The following must be done manually:
	# In gimp open fitplot.xcf and use iconify plugin to create windows object file.
	cd $SOURCE_DIR
    echo -e '\nbuilding:'$os, $arch
	go clean
	if [ "$os" == "windows" ] ; then
		# env GOOS=$1 GOARCH=$2 go build -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -H windowsgui" -v
		env CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=x86_64-w64-mingw32-gcc go install -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD` -H windowsgui" -v
	fi
	if [ "$os" == "linux" ] ; then
                # env CGO_ENABLED=1 GOOS=$1 GOARCH=$2 go install -ldflags "-X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" -v
        # Use the musl C library to statically link in libc
		env CGO_ENABLED=1 GOOS=$1 GOARCH=$2 CC=/usr/bin/musl-gcc go install -ldflags "-linkmode external -extldflags "-static"  -X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.Githash=`git rev-parse HEAD`" -v
	fi
}

function package_windows {
    echo -e 'Packaging for Windows'
	cd $PKG_DIR
	# make a backup of previous package
	sudo cp -r windows_dist/ windows_dist_old
	sudo rm -rf windows_dist
	sudo mkdir windows_dist
        cd windows_dist
        # Have replaced nw.exe icon using resource hacker.  Copy into the package.
        # Could just unzip the official version if the default icon's okay.
        # sudo unzip -q $NWJS_DIR/nwjs-v0.36.0-win-x64.zip
        sudo cp -r $NWJS_DIR/nwjs-v0.36.0-win-x64/ fitplot
	cd fitplot
	sudo mkdir nw.package
	sudo cp $BIN_DIR/windows_amd64/fitplot.exe ./nw.package/
	sudo cp $SOURCE_DIR/package.json ./nw.package/
	sudo cp $SOURCE_DIR/main.js ./nw.package/
	sudo cp $SOURCE_DIR/LICENSE.txt ./nw.package/
	sudo cp -r $SOURCE_DIR/static/ ./nw.package/
	sudo cp -r $SOURCE_DIR/tmpl/ ./nw.package/
	sudo cp -r $SOURCE_DIR/samples/ ./nw.package/
	sudo cp -r $SOURCE_DIR/db/ ./nw.package/
	sudo cp -r $SOURCE_DIR/export/ ./nw.package/
        cd ..
        sudo cp $SOURCE_DIR/'Fitplot Windows x64 Setup.nsi' . 
        sudo makensis Fitplot\ Windows\ x64\ Setup.nsi
        sudo rm Fitplot\ Windows\ x64\ Setup.nsi
        sudo rm -rf ./fitplot
        sudo sh -c "md5sum Fitplot\ Windows\ x64\ Setup.exe > md_windows.txt"
	cd $PKG_DIR
}

function package_linux {
    echo -e 'Packaging for Linux'
	cd $PKG_DIR
	# make a backup of previous package
	sudo cp -r linux_dist/ linux_dist_old
	sudo rm -rf linux_dist
	sudo mkdir linux_dist
        cd linux_dist
        sudo tar -oxzf $NWJS_DIR/nwjs-v0.36.0-linux-x64.tar.gz
        sudo mv nwjs-v0.36.0-linux-x64 fitplot
	cd fitplot
        sudo cp $SOURCE_DIR/setup.sh .
	sudo mkdir nw.package
	sudo cp $BIN_DIR/fitplot ./nw.package/
	sudo cp $SOURCE_DIR/package.json ./nw.package/
	sudo cp $SOURCE_DIR/main.js ./nw.package/
	sudo cp $SOURCE_DIR/fitplot.desktop ./nw.package/
	sudo cp -r $SOURCE_DIR/static/ ./nw.package/
	sudo cp -r $SOURCE_DIR/tmpl/ ./nw.package/
	sudo cp -r $SOURCE_DIR/samples/ ./nw.package/
	sudo cp -r $SOURCE_DIR/db/ ./nw.package/
	sudo cp -r $SOURCE_DIR/icons/ ./nw.package/
	sudo cp -r $SOURCE_DIR/export/ ./nw.package/
        cd ..
        sudo makeself fitplot fitplot.run "Fitplot by Craig Prevallet" ./setup.sh
        sudo rm -r ./fitplot/
        sudo sh -c "md5sum fitplot.run > md_linux.txt"
	cd $PKG_DIR
        # Install using: sudo ./fitplot.run --target /opt/fitplot

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
    if [ "$1" = 'all' ]; then
		build "windows" "amd64"
        package_windows
		build "linux" "amd64"
        package_linux
    fi
}



# Main entry point
# First arg is help or no args passed.
if [ "$1" = 'help' ] || [ $# -eq 0 ]; then
    echo -e '\nUsage: package {help|windows|linux|all}'
else
    date "+%nPackage started on date: %m-%d-%Y at %H:%M:%S" 
    until [ -z "$1" ]  # Until all parameters used up . . .
    do
      build_and_package $1
      shift
    done
    date "+%nPackage ended on date: %m-%d-%Y at %H:%M:%S"
fi

