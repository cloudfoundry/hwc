package hwcconfig_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/hwc/hwcconfig"
)

var _ bool = Describe("ApplicationHostConfig", func() {
	var (
		workingDirectory string
	)

	var createAllFiles = func(targetFilePath string) {
		var err error
		err = os.MkdirAll(filepath.Dir(targetFilePath), 0777)
		Expect(err).ToNot(HaveOccurred())
		err = ioutil.WriteFile(targetFilePath, []byte(""), 0666)
		Expect(err).ToNot(HaveOccurred())
	}

	var basicDeps = func(workingDirectory string) (listenPort int, rootPath string, tmpPath string, contextPath string, uuid string) {
		listenPort = 8080
		rootPath = workingDirectory + "/rootPath"
		tmpPath = workingDirectory + "/tmpPath"
		contextPath = workingDirectory + "/contextPath"
		uuid = "someuid12345"

		return
	}

	BeforeEach(func() {
		var err error

		//use the test-friendly TempDir
		workingDirectory, err = ioutil.TempDir("", "hwcconfig_test")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		_ = os.RemoveAll(workingDirectory)
	})

	Describe("Generate config file", func() {
		Context("Default config file", func() {
			It("creates default config file", func() {
				var err error

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				_, err = os.Stat(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Config file with custom modules", func() {
			var (
				modulesDirectory string
			)

			BeforeEach(func() {
				modulesDirectory = filepath.Join(workingDirectory, "modules", "hwc", "native-modules")
			})

			AfterEach(func() {
				Expect(os.Unsetenv("HWC_NATIVE_MODULES")).To(Succeed())
			})

			It("adds dlls and symlinks to applicationHost.config", func() {
				var err error

				err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				dllFilePath := filepath.Join(modulesDirectory, "exampleModule", "mymodule.dll")
				createAllFiles(dllFilePath)

				symLinkSource := filepath.Join(workingDirectory, "sourceModule.dll")
				createAllFiles(symLinkSource)

				linkFilePath := filepath.Join(modulesDirectory, "myLinkedModule", "linkModule.dll")
				err = os.MkdirAll(filepath.Dir(linkFilePath), 0777)
				Expect(err).ToNot(HaveOccurred())

				err = os.Symlink(symLinkSource, linkFilePath)
				Expect(err).ToNot(HaveOccurred())

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).ToNot(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + dllFilePath + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"myLinkedModule\" image=\"" + linkFilePath + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"myLinkedModule\" lockItem=\"true\" />"))
			})

			It("returns error when user provided directory is empty", func() {
				var err error
				emptyModulesDirectory := filepath.Join(workingDirectory, "modules")

				err = os.MkdirAll(emptyModulesDirectory, 0777)
				Expect(err).ToNot(HaveOccurred())

				err = os.Setenv("HWC_NATIVE_MODULES", emptyModulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)
				err, _ = hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("HWC_NATIVE_MODULES does not match required directory structure. See hwc README for detailed instructions."))
			})
		})
	})
})
