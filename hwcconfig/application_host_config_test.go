//go:build windows
// +build windows

package hwcconfig_test

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/hwc/hwcconfig"
)

type Configuration struct {
	XMLName         xml.Name `xml:"configuration"`
	Junk            string   `xml:"junk,attr"`
	SystemWebServer struct {
		Security struct {
			Authentication struct {
				WindowsAuthentication struct {
					Enabled                  string `xml:"enabled,attr"`
					AuthPersistNonNTLM       string `xml:"authPersistNonNTLM,attr"`
					AuthPersistSingleRequest string `xml:"authPersistSingleRequest,attr"`
					Providers                struct {
						Add []struct {
							Value string `xml:"value,attr"`
						} `xml:"add"`
					} `xml:"providers"`
				} `xml:"windowsAuthentication"`
			} `xml:"authentication"`
		} `xml:"security"`
	} `xml:"system.webServer"`
}

var _ bool = Describe("ApplicationHostConfig", func() {
	var (
		workingDirectoryPath string
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

		workingDirectoryPath, err = ioutil.TempDir("", "hwcconfig_test")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		_ = os.RemoveAll(workingDirectoryPath)
	})

	Context("With default params", func() {
		It("creates default config file", func() {
			var err error

			listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectoryPath)

			err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
			_, err = os.Stat(hwcConfig.ApplicationHostConfigPath)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("When custom modules are specified", func() {
		var (
			modulesDirectoryPath string
		)

		BeforeEach(func() {
			modulesDirectoryPath = filepath.Join(workingDirectoryPath, "modules", "hwc", "native-modules")
		})

		AfterEach(func() {
			Expect(os.Unsetenv("HWC_NATIVE_MODULES")).To(Succeed())
		})

		It("adds regular and linked DLLs to applicationHost.config", func() {
			var err error

			err = os.Setenv("HWC_NATIVE_MODULES", modulesDirectoryPath)
			Expect(err).ToNot(HaveOccurred())

			dllFilePath := filepath.Join(modulesDirectoryPath, "exampleModule", "mymodule.dll")
			createAllFiles(dllFilePath)

			linkSourcePath := filepath.Join(workingDirectoryPath, "sourceModule.dll")
			createAllFiles(linkSourcePath)

			linkFilePath := filepath.Join(modulesDirectoryPath, "myLinkedModule", "linkModule.dll")
			err = os.MkdirAll(filepath.Dir(linkFilePath), 0777)
			Expect(err).ToNot(HaveOccurred())

			err = os.Symlink(linkSourcePath, linkFilePath)
			Expect(err).ToNot(HaveOccurred())

			listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectoryPath)

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
			emptyModulesDirectoryPath := filepath.Join(workingDirectoryPath, "modules")

			err = os.MkdirAll(emptyModulesDirectoryPath, 0777)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("HWC_NATIVE_MODULES", emptyModulesDirectoryPath)
			Expect(err).ToNot(HaveOccurred())

			listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectoryPath)
			err, _ = hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("HWC_NATIVE_MODULES does not match required directory structure. See hwc README for detailed instructions."))
		})
	})

	Context("When windowsauthentication is defined", func() {
		It("adds secure windows auth config values to applicationHost.config", func() {
			var err error

			listenPort, rootPath, tmpPath, contextPath, uuid := basicDeps(workingDirectoryPath)

			err, hwcConfig := hwcconfig.New(listenPort, rootPath, tmpPath, contextPath, uuid)
			Expect(err).ToNot(HaveOccurred())
			configFileContents, err := ioutil.ReadFile(hwcConfig.ApplicationHostConfigPath)
			Expect(err).ToNot(HaveOccurred())

			var config Configuration
			err = xml.Unmarshal(configFileContents, &config)
			Expect(err).ToNot(HaveOccurred())

			Expect(config.SystemWebServer.Security.Authentication.WindowsAuthentication.Enabled).To(Equal("true"), "Missing Enabled")
			Expect(config.SystemWebServer.Security.Authentication.WindowsAuthentication.AuthPersistNonNTLM).To(Equal("true"), "Missing AuthPersistNonNTLM")
			Expect(config.SystemWebServer.Security.Authentication.WindowsAuthentication.AuthPersistSingleRequest).To(Equal("true"), "Missing AuthPersistSingleRequest")

			Expect(config.SystemWebServer.Security.Authentication.WindowsAuthentication.Providers.Add).To(HaveLen(1))
			Expect(config.SystemWebServer.Security.Authentication.WindowsAuthentication.Providers.Add[0].Value).To(Equal("Negotiate"), "Not Negotiate")
		})
	})
})
