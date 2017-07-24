package hwcconfig

import "strings"

// HwcApplication represents an application element under the IronFoundry site
type HwcApplication struct {
	PhysicalPath string
	Path         string
}

// NewHwcApplications returns the set of HwcApplications that need to be created in
// the applicationHost.config to support nested virtual directory paths. Each
// contextPath segment needs it's own application element in the applicationHost.config
func NewHwcApplications(defaultRootPath, rootPath, contextPath string) []*HwcApplication {
	var apps []*HwcApplication
	curContextPath := contextPath
	for {
		apps = append(apps, &HwcApplication{
			PhysicalPath: defaultRootPath,
			Path:         curContextPath,
		})
		nextContextPath := removeLastSegmentFromPath(curContextPath)
		if nextContextPath == curContextPath {
			break
		}
		curContextPath = nextContextPath
	}

	// only the deepest leaf context path points to the application files
	apps[0].PhysicalPath = rootPath
	return apps
}

// Removes the last segment from the path, but always returns a leading '/'
func removeLastSegmentFromPath(path string) string {
	i := strings.LastIndexByte(path, '/')
	if i > 0 {
		return path[:i]
	}
	return "/"
}
