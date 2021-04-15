$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

Install-WindowsFeature Web-WHC, Web-Webserver, Web-WebSockets, Web-ASP, Web-ASP-Net45

# Write-Host "Installing Ginkgo"
# go.exe get github.com/onsi/ginkgo/ginkgo
# if ($LastExitCode -ne 0) {
#     throw "Ginkgo installation process returned error code: $LastExitCode"
# }

go test -v -count=1 -parallel=1 ./...
if ($LastExitCode -ne 0) {
    throw "Testing hwc returned error code: $LastExitCode"
}