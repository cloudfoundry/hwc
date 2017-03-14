## hwc (hostable web core)

`hwc` is a wrapper around [Hosted Web Core API](https://msdn.microsoft.com/en-us/library/ms693832(v=vs.90).aspx) for running .NET Applications on Windows.

## Dependencies
- [Golang Windows](https://golang.org/dl/)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- Hostable Web Core
  - Install in Powershell by running `Install-WindowsFeature Web-WHC`

### Building

```
./scripts/build.sh

```

### Test

Unit Tests:

```
ginkgo -r -race
```

### Running

When web applications are pushed to Cloud Foundry they are pushed out to one or more Windows cells and run via `hwc.exe`. For development purposes you can run an ASP.NET web application much like IISExpress by directly invoking `hwc.exe`.

1. Install the following Windows features: Hostable Web Core, ASP.NET 4.6, Websockets.
1. [Build](https://github.com/cloudfoundry-incubator/hwc#building) hwc.exe, or [Download](https://github.com/cloudfoundry-incubator/hwc/releases/) the prebuilt binary from the GitHub releases page.
1. From PowerShell start the web server: `& { $env:PORT=8080; .\hwc.exe -appRootPath "C:\wwwroot\inetpub\myapproot" }`. Ensure the appRootPath points to a directory with a ready to run ASP.NET application.

You should now be able to browse to `http://localhost:8080/` and even attach a debugger and set breakpoints to the `hwc.exe` process if so desired.
