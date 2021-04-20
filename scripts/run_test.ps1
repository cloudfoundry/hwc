$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

# Install windows features
Install-WindowsFeature Web-WHC, Web-Webserver, Web-WebSockets, Web-ASP, Web-ASP-Net45

# Install chocolatey
Invoke-WebRequest https://chocolatey.org/install.ps1 -UseBasicParsing | Invoke-Expression

# Install go and gcc
choco install -y mingw

# Run tests
go test -v -count=1 -parallel=1 ./...
if ($LastExitCode -ne 0) {
    throw "Testing hwc returned error code: $LastExitCode"
}