$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

Install-WindowsFeature Web-WHC
Install-WindowsFeature Web-Webserver
Install-WindowsFeature Web-WebSockets
Install-WindowsFeature Web-WHC
Install-WindowsFeature Web-ASP
Install-WindowsFeature Web-ASP-Net45

cd hwc

Write-Host "Installing Ginkgo"
go.exe get github.com/onsi/ginkgo/ginkgo
if ($LastExitCode -ne 0) {
    throw "Ginkgo installation process returned error code: $LastExitCode"
}

~/go/bin/ginkgo.exe -r -race -keepGoing -p
if ($LastExitCode -ne 0) {
    throw "Testing hwc returned error code: $LastExitCode"
}