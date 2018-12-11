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

	Describe("Generate config file", func() {
		Context("Default config file", func() {
			It("creates default config file", func() {
				var tempDir string
				var rootPath string
				var tmpPath string
				var contextPath string
				var err error
				tempDir, err = ioutil.TempDir("", "hwcconfig_test")
				Expect(err).NotTo(HaveOccurred())

				rootPath = tempDir + "/rootPath"
				tmpPath = tempDir + "/tmpPath"
				contextPath = tempDir + "/contextPath"
				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				_, err = os.Stat(hwcConfig.ApplicationHostConfigPath)
				Expect(err).NotTo(HaveOccurred())

				os.RemoveAll(tempDir)
			})
		})

		Context("Config file with custom modules", func() {
			var tempDir string

			AfterEach(func() {
				Expect(os.Unsetenv("HWC_NATIVE_MODULES")).To(Succeed())
				_ = os.RemoveAll(tempDir)
			})

			It("creates config file with a module", func() {
				var rootPath string
				var tmpPath string
				var contextPath string
				var modulesRoot string
				var err error

				tempDir, err = ioutil.TempDir("", "hwcconfig_test_single")
				Expect(err).NotTo(HaveOccurred())
				rootPath = tempDir + "/rootPath"
				tmpPath = tempDir + "/tmpPath"
				contextPath = tempDir + "/contextPath"
				modulesRoot = tempDir + "/modules"

				modulesDir := filepath.Join(modulesRoot, "hwc", "native-modules")
				pathToDll := filepath.Join(modulesDir, "exampleModule", "mymodule.dll")
				err = os.Setenv("HWC_NATIVE_MODULES", modulesDir)
				Expect(err).NotTo(HaveOccurred())

				writeDllToPath(pathToDll)

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(err).NotTo(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + pathToDll + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))
			})

			It("creates a config file with multiple modules", func() {
				var rootPath string
				var tmpPath string
				var contextPath string
				var modulesRoot string
				var err error

				tempDir, err = ioutil.TempDir("", "hwcconfig_test_multi")
				Expect(err).NotTo(HaveOccurred())

				rootPath = filepath.Join(tempDir, "rootPath")
				tmpPath = filepath.Join(tempDir, "tmpPath")
				contextPath = filepath.Join(tempDir, "contextPath")
				modulesRoot = filepath.Join(tempDir, "modules")

				modulesDir := filepath.Join(modulesRoot, "hwc", "native-modules")
				//fmt.Fprintf(os.Stderr, "\nERR: expected %s\n", modulesDir)
				//fmt.Fprintf(os.Stderr, "ERR: before   %s\n", os.Getenv("HWC_NATIVE_MODULES"))
				err = os.Setenv("HWC_NATIVE_MODULES", modulesDir)
				//fmt.Fprintf(os.Stderr, "ERR: after    %s\n", os.Getenv("HWC_NATIVE_MODULES"))
				//fmt.Fprintf(os.Stderr, "ERR: tmpPath  %s\n", tmpPath)
				Expect(err).NotTo(HaveOccurred())

				pathToDll1 := filepath.Join(modulesDir, "exampleModule", "mymodule.dll")
				writeDllToPath(pathToDll1)

				pathToDll2 := filepath.Join(modulesDir, "anotherModule", "anotherModule.dll")
				writeDllToPath(pathToDll2)

				err, hwcConfig := hwcconfig.New(8080, rootPath, tmpPath, contextPath, "someuid12345")
				Expect(err).NotTo(HaveOccurred())
				configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" image=\"" + pathToDll1 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"exampleModule\" lockItem=\"true\" />"))

				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" image=\"" + pathToDll2 + "\""))
				Expect(string(configFileContents)).To(ContainSubstring("<add name=\"anotherModule\" lockItem=\"true\" />"))
			})
			//Ensure symlinks work
			//Sort order of DLL files
			//Multiple modules included
			//Dont add empty directories
		})
	})
})

func writeDllToPath(pathToDll1 string) {
	var err error

	err = os.MkdirAll(filepath.Dir(pathToDll1), 0777)
	Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(pathToDll1, []byte(""), 0666)
	Expect(err).NotTo(HaveOccurred())
}
