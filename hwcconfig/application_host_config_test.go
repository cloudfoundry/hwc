package hwcconfig_test

import (
	"log"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/hwc/hwcconfig"
)

var _ = Describe("ApplicationHostConfig", func() {
	var tempDir string
	var rootPath string
	var tmpPath string
	var contextPath string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir(os.TempDir(), "hwcconfig_test")
		if err != nil {
			log.Fatal(err)
		}

		rootPath = tempDir + "/rootPath"
		tmpPath = tempDir + "/tmpPath"
		contextPath = tempDir + "/contextPath"
	})
	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Fatal(err)
		}
	})

	Describe("Generate config file", func() {
		Context("Default config file", func() {
			err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
			_, err = os.Stat(hwcConfig.ApplicationHostConfigPath)
			It("creates default config file", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Config file with custom modules", func() {
			fakeDll := tempDir + "someModule.dll"
			It("creates fake dll file", func() {
				_, err := os.Create(fakeDll)
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets CUSTOMMODULES environment variable", func() {
				err := os.Setenv("CUSTOMMODULES", "someModule," + fakeDll)
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates config file with custom modules", func() {
				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(err).NotTo(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"someModule\" image=\"" + fakeDll + "\""))
			})
		})
	})
})
