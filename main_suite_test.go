package main_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var hwcBinPath string

func TestHWC(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(200 * time.Millisecond)
	RunSpecs(t, "HWC")
}

var randomGenerator *rand.Rand

var _ = BeforeSuite(func() {
	var err error

	randomGenerator = rand.New(rand.NewSource(GinkgoRandomSeed() + int64(GinkgoParallelProcess())))

	hwcBinPath, err = gexec.BuildWithEnvironment("code.cloudfoundry.org/hwc", []string{"CGO_ENABLED=1", "GO_EXTLINK_ENABLED=1"})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
