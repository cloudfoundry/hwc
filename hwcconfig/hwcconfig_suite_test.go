package hwcconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHwcconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hwcconfig Suite")
}
