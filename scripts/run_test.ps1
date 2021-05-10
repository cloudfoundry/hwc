$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

Install-WindowsFeature Web-WHC
Install-WindowsFeature Web-Webserver
Install-WindowsFeature Web-WebSockets
Install-WindowsFeature Web-WHC
Install-WindowsFeature Web-ASP
Install-WindowsFeature Web-ASP-Net45

Write-Host "Installing Ginkgo"
go.exe get github.com/onsi/ginkgo/ginkgo
if ($LastExitCode -ne 0) {
    throw "Ginkgo installation process returned error code: $LastExitCode"
}

~/go/bin/ginkgo.exe -r -race -keepGoing -p
if ($LastExitCode -ne 0) {
    throw "Testing hwc returned error code: $LastExitCode"
}

################
#
# $ErrorActionPreference = "Stop";
# trap { $host.SetShouldExit(1) }
#
# # Install windows features
# Install-WindowsFeature Web-WHC, Web-Webserver, Web-WebSockets, Web-ASP, Web-ASP-Net45
#
# # Install chocolatey
# Invoke-WebRequest https://chocolatey.org/install.ps1 -UseBasicParsing | Invoke-Expression
#
# # Install go and gcc
# choco install -y mingw
#
# # Run tests
# go test -v -count=1 -parallel=1 ./...
# if ($LastExitCode -ne 0) {
#     throw "Testing hwc returned error code: $LastExitCode"
# }
