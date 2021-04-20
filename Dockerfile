FROM golang:1.16.3-windowsservercore-1809
#FROM mcr.microsoft.com/windows/servercore:ltsc2019

WORKDIR /tests
COPY . .
CMD ["powershell", "C:/tests/scripts/run_test.ps1"]