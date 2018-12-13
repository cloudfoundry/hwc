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
	// 	tmpPath          string
	// 	rootPath         string
	// 	contextPath      string
	// )

	var writeDllToPath = func(pathToDll1 string) {
		var err error
		err = os.MkdirAll(filepath.Dir(pathToDll1), 0777)
		Expect(err).ToNot(HaveOccurred())
		err = ioutil.WriteFile(pathToDll1, []byte(""), 0666)
		Expect(err).ToNot(HaveOccurred())
	}

	var basicDeps = func(workingDirectory string) (listenPort int, rootPath string, tmpPath string, contextPath string, uuid string) {
		listenPort = 8080
		rootPath = workingDirectory + "/rootPath"
		tmpPath = workingDirectory + "/tmpPath"
		contextPath = workingDirectory + "/contextPath"
		uuid = "someuid12345"

		// return "", rootPath, tmpPath, contextPath
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

			It("creates config file with a module", func() {
				var err error

				pathToDll := filepath.Join(modulesDirectory, "exampleModule", "mymodule.dll")
				err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				writeDllToPath(pathToDll)

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).ToNot(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + pathToDll + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))
			})

			It("creates a config file with multiple modules", func() {
				var err error

				err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				pathToDll1 := filepath.Join(modulesDirectory, "exampleModule", "mymodule.dll")
				writeDllToPath(pathToDll1)

				pathToDll2 := filepath.Join(modulesDirectory, "anotherModule", "anotherModule.dll")
				writeDllToPath(pathToDll2)

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).ToNot(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + pathToDll1 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" image=\"" + pathToDll2 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" lockItem=\"true\" />"))
			})

			It("returns error when user defined directory does NOT exist", func() {
				err := os.Setenv("HWC_NATIVE_MODULES", filepath.Join("some", "nonexistent", "path"))
				Expect(err).ToNot(HaveOccurred())

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				hwcConfigErr, _ := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(hwcConfigErr).To(HaveOccurred())
				Expect(hwcConfigErr).To(MatchError("Path \"some\\nonexistent\\path\" does not exist"))
			})

			It("ignores empty sub-directories in user defined path", func() {
				var err error

				err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				emptyModuleDirectoryPath := filepath.Join(modulesDirectory, "emptyDirectory")

				err = os.MkdirAll(emptyModuleDirectoryPath, 0777)
				Expect(err).ToNot(HaveOccurred())

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).ToNot(HaveOccurred())

				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).ToNot(ContainSubstring("emptyDirectory"))
			})

			It("appends symlinked modules to application host", func() {
				var err error

				err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(err).ToNot(HaveOccurred())

				moduleDirectory := filepath.Join(modulesDirectory, "myLinkedModule")
				err = os.MkdirAll(moduleDirectory, 0777)
				Expect(err).ToNot(HaveOccurred())

				symLinkSource := filepath.Join(workingDirectory, "sourceModule.dll")
				writeDllToPath(symLinkSource)
				linkedModule := filepath.Join(moduleDirectory, "linkModule.dll")
				err = os.Symlink(symLinkSource, linkedModule)
				Expect(err).ToNot(HaveOccurred())

				listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectory)

				err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
				Expect(err).ToNot(HaveOccurred())

				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"myLinkedModule\" image=\"" + linkedModule + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"myLinkedModule\" lockItem=\"true\" />"))
			})
		})
	})
})
