<img src="https://github.com/cprevallet/fitplot/blob/master/icons/fitplot_color.png">

# fitplot
Visualize, summarize and analyze running activities from Garmin devices.

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
                go get github.com/mitchellh/go-homedir
                go get github.com/mattn/go-sqlite3
		go get bitbucket.org/liamstask/goose/cmd/goose
		go get github.com/jezard/fit
		go get github.com/cprevallet/fitplot
	2. Build:
		go get github.com/cprevallet/fitplot
		go install fitplot
	3. Package (optional):
		Install dependencies for package.sh (see comments in package.sh for list)
		/.package.sh {help|windows|linux|all}

Prebuilt Binaries
-----------------

# Installation

Before you start using Fitplot, you have to make it available on your
computer. Even if it’s already installed, it’s probably a good idea to
update to the latest version. You can either install it as a package or
via another installer, or download the source code and compile it
yourself.

> **Important**
> 
> There are a number of prerequisites in order to use this application.
> 
>   - Builds are available for 64-bit versions of Microsoft Windows (7
>     and above) as well as 64-bit Linux distributions. MacOS/OSX builds
>     are not available or supported.
> 
>   - A 64-bit Intel-compatible processor.
> 
>   - Administrator/root privileges are necessary for installation.
> 
>   - A functioning Internet connection is necessary to use the
>     application.
> 
>   - A method (USB or ANT) to transfer .FIT or .TCX files to the PC.

The latest release builds are available for download here:

<https://github.com/cprevallet/fitplot/releases>

Fitplot is open source software. The actual program code for this
software to view and modify is online at Github.

<https://github.com/cprevallet/fitplot>

## Installing on Windows(64 bit)

Installation is a performed simply by downloading and running 'Fitplot
Windows x64 Setup.exe' as Administrator. The setup executable will
display the license and then prompt for an installation location.

A start menu folder will be created containing links to start and to
uninstall Fitplot.

## Installing on GNU/Linux(64 bit)

Installation is performed by first setting the executable bit on the
script and the running it (as root) to copy the files into the
/opt/fitplot directory on the user’s drive.

A desktop application menu link will be created (under Accessories) to
start Fitplot.

``` console
$ sudo chmod +x fitplot.run
$ sudo ./fitplot.run --target /opt/fitplot
Creating directory /opt/fitplot
Verifying archive integrity...  100%   All good.
Uncompressing Fitplot by Craig Prevallet  100%
Fitplot has extracted itself.
Setting up desktop and icons files.
All done.
```

> **Tip**
> 
> How do I know this software hasn’t been modified from the original
> source?
> 
> In lieu of a "signing certificate" that both Apple and Microsoft
> support within their respective operating systems, Fitplot is offered
> with a "message digest" to insure that the application has not been
> tampered with. As a result, when installing under Windows, a warning
> about an unknown publisher may be issued.
> 
> The following utilities may be executed from the command line after
> downloading the builds to verify the provided message digests.
> Performing this check is optional but a good practice.
> 
>   - Windows
> 
> 
> 
> ``` console
> C:\ CertUtil -hashfile 'Fitplot Windows x64 Setup.exe' MD5
> ```
> 
>   - Linux
> 
>
> 
> ``` console
> $ md5sum fitplot.run
> ```
> 
> The resulting displayed checksum should match the checksum value in
> the file md\_windows.txt or md\_linux.txt.


Additional Help and Support
---------------------------

The application has an integrated help system which may be accessed from
the menu once the application is started.

The tracking system at Github will be used to report problems and suggest enhancements.
https://github.com/cprevallet/fitplot/issues
                                                                                                        
License
-------

This software is governed by the following software license:

    Copyright 2016 - 2019 Craig S. Prevallet

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.


