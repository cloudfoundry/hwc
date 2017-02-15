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
