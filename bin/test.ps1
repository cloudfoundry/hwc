$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-WebServer
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-WebSockets
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-HostableWebCore
Enable-WindowsOptionalFeature -Online -All -FeatureName IIS-ASPNET45

Invoke-Expression "go run github.com/onsi/ginkgo/v2/ginkgo $args"
if ($LastExitCode -ne 0) {
  throw "tests failed"
}
