package contextpath

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

func Default() string {
	return "/"
}

func New(appEnv *cfenv.App) (string, error) {
	return appContextPath(appEnv)
}

func appContextPath(appEnv *cfenv.App) (string, error) {
	contextPath := Default()
	uniqueContextPaths := map[string]bool{}
	for _, applicationURI := range appEnv.ApplicationURIs {
		contextPath = parseContextPath(applicationURI)
		uniqueContextPaths[contextPath] = true
	}
	err := checkContextPathIsUnique(uniqueContextPaths)
	if err != nil {
		return "", err
	}
	return contextPath, nil
}

func parseContextPath(applicationURI string) string {
	parts := strings.Split(applicationURI, "/")
	return "/" + strings.TrimSuffix(strings.Join(parts[1:], "/"), "/")
}

func checkContextPathIsUnique(uniqueContextPaths map[string]bool) error {
	if len(uniqueContextPaths) <= 1 {
		return nil
	}
	var errParts []string
	for contextPath := range uniqueContextPaths {
		errParts = append(errParts, contextPath)
	}
	sort.Strings(errParts)
	return fmt.Errorf("Application may not contain conflicting route paths: %s", strings.Join(errParts, ", "))
}
