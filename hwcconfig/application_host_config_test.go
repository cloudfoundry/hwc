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
		tmpPath          string
		rootPath         string
		contextPath      string
	)

	var writeDllToPath = func(pathToDll1 string) {
		var err error
		err = os.MkdirAll(filepath.Dir(pathToDll1), 0777)
		Expect(err).ToNot(HaveOccurred())
		err = ioutil.WriteFile(pathToDll1, []byte(""), 0666)
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		var err error
		workingDirectory, err = ioutil.TempDir("", "hwcconfig_test")
		Expect(err).ToNot(HaveOccurred())

		rootPath = workingDirectory + "/rootPath"
		tmpPath = workingDirectory + "/tmpPath"
		contextPath = workingDirectory + "/contextPath"
	})

	AfterEach(func() {
		_ = os.RemoveAll(workingDirectory)
	})

	Describe("Generate config file", func() {
		Context("Default config file", func() {
			It("creates default config file", func() {
				var err error

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
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

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
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

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(err).ToNot(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + pathToDll1 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" image=\"" + pathToDll2 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" lockItem=\"true\" />"))
			})

			It("returns error when user defined directory does NOT exist", func() {
				envErr := os.Setenv("HWC_NATIVE_MODULES", filepath.Join("some", "nonexistent", "path"))
				Expect(envErr).ToNot(HaveOccurred())

				hwcConfigErr, _ := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(hwcConfigErr).To(HaveOccurred())
				Expect(hwcConfigErr).To(MatchError("Path \"some\\nonexistent\\path\" does not exist"))
			})

			It("ignores empty sub-directories in user defined path", func() {
				envErr := os.Setenv("HWC_NATIVE_MODULES", modulesDirectory)
				Expect(envErr).ToNot(HaveOccurred())

				emptySubDirectoryPath := filepath.Join(modulesDirectory, "emptyDirectory")

				directoryErr := os.MkdirAll(emptySubDirectoryPath, 0777)
				Expect(directoryErr).ToNot(HaveOccurred())

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(err).ToNot(HaveOccurred())

				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(configFileContents)).ToNot(ContainSubstring("emptyDirectory"))
			})
		})
	})
})
