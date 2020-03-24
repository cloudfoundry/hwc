if (-not (Test-Path env:GOPATH)) {
  Write-Output 'GOPATH environment variable must be set before running build!'
  exit 1
}

New-Item -path "$PWD/hwc-rel" -type directory -force

$env:CGO_ENABLED=1
$env:GO_EXTLINK_ENABLED=1
# x64
$env:CC="x86_64-w64-mingw32-gcc"
$env:GOARCH="amd64"
go build -o $PWD/hwc-rel/hwc.exe code.cloudfoundry.org/hwc

# Win32
$env:CC="i686-w64-mingw32-gcc"
$env:GOARCH="386"
go build -o $PWD/hwc-rel/hwc_x86.exe code.cloudfoundry.org/hwc
