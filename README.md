<img src="https://github.com/cprevallet/fitplot/blob/master/icons/fitplot_color.png">

# fitplot
Visualize, summarize and analyze runnning activities from Garmin devices.


Application Requirements
------------------------

### Operating System

Builds will be made available for Windows(64-bit), Linux (64 bit) and Mac 
OSX (64 bit).

### Input devices

A mouse or trackpad rather than a touchscreen is recommended in order to
pan and scroll the graph generated by the application. A touchscreen may
be used other features.

### Browser

Fitplot requires access to the internet and a web browser.

#### Browser version

In theory, any browser capable of supporting a combination of Plot.ly,
Material Design Lite components, and Google Maps Javascript API and
Google Charts should be capable of using the application.

In general browsers with the following version numbers or later should
be compatible (but see the tested configurations section below):

-   Microsoft Internet Explorer 10 or above
-   Mozilla Firefox version 31 or above
-   Google Chrome version 31 or above
-   Apple Safari version 9.1 or above

#### Browser settings

Javascript and cookies must be enabled to use the application. This is
generally the out-of-the-box setting for most modern browsers.

### Tested configurations

Fitplot has been tested under the following combinations of operating
system and browser:

-   Microsoft Windows 7 (64 bit build)
    :   -   Internet Explorer 11
        -   Firefox 47.0.1
        -   Chrome 51.0.2704.103

-   Debian GNU/Linux 8 (Jessie - 64 bit)
    :   -   Chrome 51.0.2704.106
        -   Firefox 45.2.0

-   Apple OS X El Capitan version 10.11.6
    :   -   Safari 9.1.2

### Input file formats

Fitplot will display files generated in either Garmin's FIT or TCX
formats.

Digital Security
----------------

### Application Integrity

In lieu of a "signing certificate" that both Apple and Microsoft support
within their respective operating systems, Fitplot is offered with a
"message digest" to insure that the application has not been tampered
with.

Background information on message digests and digital certificates (in
the context of another fine program known as "Pretty Good Privacy") may
be found here:

-   <http://www.pgpi.org/doc/pgpintro/#p12>
-   <http://www.pgpi.org/doc/pgpintro/#p14>

The following utilities may be executed from the command line to verify
the provided message digests upon download.

-   Windows

<!-- -->

    CertUtil -hashfile [file.ext] MD5

-   Mac OSX

<!-- -->

    md5 -r [file.ext]

-   Linux

<!-- -->

    md5sum [file.ext]

If you are concerned about security, run these utilities on the
downloaded application files and be sure the message digest matches.

Building From Source
--------------------

### Background
Fitplot source code is available on GitHub. 

-   <https://github.com/cprevallet/fitplot>

### Dependencies
This software depends on the excellent fit package.

-  <https://github.com/jezard/fit>

### Building and Packaging
The software source code is an (unholy) mixture of Golang (version 1.6.2) and 
Javascript, HTML and Javascript.  The source code may be cross-compiled and packaged by using 
the shell script (package.sh) under Linux.  Building under Windows is 
untested. See comments in that document about the other tools that may 
be required in the build environment. 

	1. Setup:  
		Download Go from https://golang.org/dl/
		go get github.com/jezard/fit
		go get github.com/cprevallet/fitplot
	2. Build:
		go get github.com/cprevallet/fitplot
		go build fiplot
	3. Package (optional): 
		Install dependencies for package.sh (see comments in package.sh for list)
		/.package.sh {help|windows|linux|osx|all}

Installation
------------

-   Windows(64 bit)

Installation is a performed simply by downloading and running Fitplot
Windows x64 Setup.exe. The setup executable will display the license and
then prompt for an installation location. A start menu icon and folder
will be created.

-   GNU/Linux(64 bit)

Installation is performed by copying the files from the delivery medium
and installing into the /opt/fitplot directory on the user's drive.

    sudo tar -xvzf fitplot.tgz -C /opt/
    sudo /opt/fitplot/icons/cpfitplot_color.png /usr/share/icons/hicolor/128x128/apps/
    sudo cp /opt/fitplot/fitplot.desktop /usr/share/applications/

-   Mac OSX (64 bit)

Installation is performed by downloading the file with the dmg file and
single clicking on it. This should result in a drive icon appearing on
the desktop. Double click on it to open. Proceed to Starting the
Application.

The program needs write access to a temporary directory (typically
C:\\Users\\User Name\\AppData\\Local\\Temp on MS Windows) or (/tmp on
Linux and OSX). Nothing else is required.



### Starting the application

Fitplot has both a web server and web client. Both must be loaded in
order to use the application.

-   Windows
    :   -   Start Menu
        -   Fitplot
        -   You will receive a message indicating that the application
            is an unsigned binary from an unknown developer and asking
            if you are sure you want to run it. See the Digital Security
            and Privacy section.

-   Linux
    :   -   From a bash shell: /opt/fitplot/fitplot.sh
        -   From the menu (if desktop file was copied per the
            installation instructions): /Utilities/Fitplot

-   OSX
    :   -   Click or tap with two fingers on fitplot.command to open the
            application.
        -   You will receive a message indicating that fitplot.command
            is from an unknown developer and asking if you are sure you
            want to open it. This is due to the developer (me) not
            signing and making it available to the Mac App Store. See
            the Digital Security and Privacy section.
        -   Click the open button to begin the application.

A terminal window may appear and the application will start as a tab in
the user's default browser.

If the browser client is closed and the server is left running, the
user-interface may be generated by opening any supported browser and
typing "<http://localhost:8080>" (without the quotes) into the address
bar.

Additional Help and Support
---------------------------

The application has an integrated help system which may be accessed from
the menu once the application is started.

The tracking system at Github will be used to report problems and suggest enhancements.
https://github.com/cprevallet/fitplot/issues
                                                                                                        
License
-------

This software is governed by the following software license:

    Copyright 2016 Craig S. Prevallet

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

