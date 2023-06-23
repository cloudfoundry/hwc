package contextpath_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContextpath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contextpath Suite")
}
