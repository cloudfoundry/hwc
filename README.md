## hwc (hostable web core)

`hwc` is a wrapper around [Hosted Web Core API](https://msdn.microsoft.com/en-us/library/ms693832(v=vs.90).aspx) for running .NET Applications on Windows.

## Dependencies
- [Golang](https://golang.org/dl/)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [MinGW-w64](https://sourceforge.net/projects/mingw-w64/)

## Compiling on Windows 

1. Install Golang from the [golang downloads page](https://golang.org/dl/).
2. Install MinGW-w64 using the [mingw-w64-install.exe](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/). Select the `x86_64` Architecture.
3. Ensure you've set the GOPATH environment variable.
4. Ensure you've added the `x86_64-w64-mingw32-gcc` compiler to your Windows PATH. It will be installed to a directory similar to: `C:\Program Files\mingw-w64\x86_64-7.2.0-posix-seh-rt_v5-rev1\mingw64\bin`

```PowerShell
git clone --recursive git@github.com:cloudfoundry-incubator/hwc "$env:GOPATH/src/code.cloudfoundry.org/hwc"
cd "$env:GOPATH/src/code.cloudfoundry.org/hwc"
.\scripts\build.ps1
```

## Cross Compiling on OSX

1. Install Golang from the [golang downloads page](https://golang.org/dl/).
2. Install MinGW-w64, `brew install mingw-w64`
3. Ensure you've set the GOPATH environment variable.

```
git clone --recursive git@github.com:cloudfoundry-incubator/hwc $GOPATH/src/code.cloudfoundry.org/hwc
cd $GOPATH/src/code.cloudfoundry.org/hwc
./scripts/build.sh
```

## Cross Compiling on Linux

1. Install Golang `sudo apt-get install gccgo-go` 
2. Install MinGW-w64, `sudo apt-get install mingw-w64`
3. Ensure you've set the GOPATH environment variable.

```
git clone --recursive git@github.com:cloudfoundry-incubator/hwc $GOPATH/src/code.cloudfoundry.org/hwc
cd $GOPATH/src/code.cloudfoundry.org/hwc
./scripts/build.sh
```

## Running the Tests (Windows Only)

Install Windows/.NET Features from Powershell by running:
```
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-WebServer
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-WebSockets
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-HostableWebCore
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-ASPNET45
```

Install ginkgo:
```
go get github.com/onsi/ginkgo/ginkgo
```

Execute the tests:
```
& "$env:GOPATH\bin\ginkgo.exe" -r -race
```

This will run the test suite which spins up several web applications hosted under hwc.exe and validates the behavior.

## Running (Windows Only)

When web applications are pushed to Cloud Foundry they are pushed out to one or more Windows cells and run via `hwc.exe`. For development purposes you can run an ASP.NET web application much like IISExpress by directly invoking `hwc.exe`.

1. Install the Windows features in the "Running the Tests" section above.
1. Build hwc.exe, or [Download](https://github.com/cloudfoundry-incubator/hwc/releases/) the prebuilt binary from the GitHub releases page.
1. From PowerShell start the web server: `& { $env:PORT=8080; .\hwc.exe -appRootPath "C:\wwwroot\inetpub\myapproot" }`. Ensure the appRootPath points to a directory with a ready to run ASP.NET application.

You should now be able to browse to `http://localhost:8080/` and even attach a debugger and set breakpoints to the `hwc.exe` process if so desired.
