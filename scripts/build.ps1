if (-not (Test-Path env:GOPATH)) {
  Write-Output 'GOPATH environment variable must be set before running build!'
  exit 1
}

New-Item -path "$PWD/hwc-rel" -type directory -force
$env:CGO_ENABLED=1
$env:GO_EXTLINK_ENABLED=1
$env:CC="x86_64-w64-mingw32-gcc"

go build -o $PWD/hwc-rel/hwc.exe code.cloudfoundry.org/hwc
