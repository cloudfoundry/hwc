//go:build windows
// +build windows

package main_test

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("HWC", func() {
	Context("when the app PORT is not set", func() {
		It("errors", func() {
			app := startAppWithEnv("nora", []string{"PORT="}, false)
			Eventually(app.session).Should(gexec.Exit(1))
			Eventually(app.session.Err).Should(gbytes.Say("Missing PORT environment variable"))
			stopApp(app)
		})
	})

	Context("when the app USERPROFILE is not set", func() {
		It("errors", func() {
			app := startAppWithEnv("nora", []string{"USERPROFILE="}, false)
			Eventually(app.session).Should(gexec.Exit(1))
			Eventually(app.session.Err).Should(gbytes.Say("Missing USERPROFILE environment variable"))
			stopApp(app)
		})
	})

	Context("Given that I am missing a required .dll", func() {
		It("errors", func() {
			app := startAppWithEnv("nora", []string{"WINDIR="}, false)
			Eventually(app.session).Should(gexec.Exit(1))
			Eventually(app.session.Err).Should(gbytes.Say("Missing required DLLs:"))
			stopApp(app)
		})
	})

	Context("Given native module environment variable is set", func() {
		var someDir, otherDir string

		BeforeEach(func() {
			var err error

			someDir, err = os.MkdirTemp("", "some-modules-dir")
			Expect(err).NotTo(HaveOccurred())

			err = os.MkdirAll(filepath.Join(someDir, "SomeModule"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(someDir, "SomeModule", "some-module.html"), nil, 0777)
			Expect(err).ToNot(HaveOccurred())

			otherDir, err = os.MkdirTemp("", "other-modules-dir")
			Expect(err).NotTo(HaveOccurred())

			err = os.MkdirAll(filepath.Join(otherDir, "OtherModule"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(otherDir, "OtherModule", "other-module.html"), nil, 0777)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(someDir)).To(Succeed())
			Expect(os.RemoveAll(otherDir)).To(Succeed())
		})

		It("exits with an error when the module directory contains an invalid DLL", func() {
			app := startAppWithEnv("nora", []string{fmt.Sprintf("HWC_NATIVE_MODULES=%s;%s", someDir, otherDir)}, false)
			Eventually(app.session).Should(gexec.Exit(1))
			Eventually(app.session).Should(gbytes.Say("HWC loading native module: .*some-module.html"))
			Eventually(app.session).Should(gbytes.Say("HWC loading native module: .*other-module.html"))
			Eventually(app.session.Err).Should(gbytes.Say("HWC Failed to start: return code:"))
			stopApp(app)
		})
	})

	Context("Given that I have http compression", func() {
		var app hwcApp

		BeforeEach(func() {
			app = startApp("ASPNetTemplateApplication")
			Eventually(app.session).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gexec.Exit(0))
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
		})

		It("does static compression", func() {
			staticUrl := fmt.Sprintf("http://localhost:%d/Content/bootstrap.css", app.port)
			var header http.Header

			//the app needs more than one request in order to do the compression
			for i := 0; i < 2; i++ {
				header = successfulRequest(staticUrl)
			}

			Expect(header["Content-Encoding"]).To(ContainElement("gzip"))
			Expect(header["Vary"]).To(ContainElement("Accept-Encoding"))

			cachePath := filepath.Join(app.profileDir, "tmp", "IIS Temporary Compressed Files", fmt.Sprintf("AppPool%d", app.port), "$^_gzip_C^")
			cacheInfo, err := os.Stat(cachePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(cacheInfo.IsDir())
			Expect(header["Content-Length"]).To(ContainElement("22748"))
		})

		It("does dynamic compression", func() {
			dynamicUrl := fmt.Sprintf("http://localhost:%d/Content/css", app.port)
			var header http.Header

			//the app needs more than one request in order to do the compression
			for i := 0; i < 2; i++ {
				header = successfulRequest(dynamicUrl)
			}
			Expect(header["Content-Encoding"]).To(ContainElement("gzip"))
			Expect(header["Vary"]).To(ContainElement("User-Agent,Accept-Encoding"))
			Expect(header["Content-Length"]).To(ContainElement("23434"))
		})
	})

	Context("Given that I have a nora with http compression", func() {
		var app hwcApp
		BeforeEach(func() {
			app = startApp("nora")
			Eventually(app.session).Should(gbytes.Say("Server Started"))
		})
		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gexec.Exit(0))
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
		})
		It("does dynamic compression for JSON", func() {
			dynamicUrl := fmt.Sprintf("http://localhost:%d/healthcheck", app.port)
			var header http.Header
			//the app needs more than one request in order to do the compression
			for i := 0; i < 2; i++ {
				header = successfulRequest(dynamicUrl)
			}
			Expect(header["Content-Type"]).To(ContainElement("application/json; charset=utf-8"))
			Expect(header["Content-Encoding"]).To(ContainElement("gzip"))
			Expect(header["Content-Length"]).To(ContainElement("52"))
			Expect(header["Vary"]).To(ContainElement("Accept-Encoding"))
		})
		It("does static compression for image", func() {
			staticUrl := fmt.Sprintf("http://localhost:%d/Content/art.jpg", app.port)
			var header http.Header
			//the app needs more than one request in order to do the compression
			for i := 0; i < 2; i++ {
				header = successfulRequest(staticUrl)
			}
			Expect(header["Content-Type"]).To(ContainElement("image/jpeg"))
			Expect(header["Content-Encoding"]).To(ContainElement("gzip"))
			Expect(header["Content-Length"]).To(ContainElement("7630"))
			Expect(header["Vary"]).To(ContainElement("Accept-Encoding"))
		})
	})

	Context("Given that I have an ASP.NET MVC application (nora)", func() {
		var app hwcApp

		BeforeEach(func() {
			app = startApp("nora")
			Eventually(app.session).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gexec.Exit(0))
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
		})

		It("runs it on the specified port", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := io.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal(fmt.Sprintf(`"hello i am %s running on http://localhost:%d/"`, "nora", app.port)))
		})

		It("correctly utilizes the USERPROFILE temp directory", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			Expect(filepath.Join(app.profileDir, "tmp", "root")).To(BeADirectory())
			Expect(err).ToNot(HaveOccurred())

			By("placing config files in the temp directory", func() {
				Expect(filepath.Join(app.profileDir, "tmp", "config", "Web.config")).To(BeAnExistingFile())
				Expect(filepath.Join(app.profileDir, "tmp", "config", "ApplicationHost.config")).To(BeAnExistingFile())
				Expect(filepath.Join(app.profileDir, "tmp", "config", "Aspnet.config")).To(BeAnExistingFile())
			})

			By("creating an IIS Temporary Compressed Files directory", func() {
				Expect(filepath.Join(app.profileDir, "tmp", "IIS Temporary Compressed Files")).To(BeADirectory())
			})

			By("creating an ASP Compiled Templates directory", func() {
				Expect(filepath.Join(app.profileDir, "tmp", "ASP Compiled Templates")).To(BeADirectory())
			})
		})

		It("does not add unexpected custom headers", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			var customHeaders []string
			for h := range res.Header {
				if strings.HasPrefix(h, "X-") {
					customHeaders = append(customHeaders, strings.ToLower(h))
				}
			}
			Expect(len(customHeaders)).To(Equal(1))
			Expect(customHeaders[0]).To(Equal("x-aspnet-version"))
		})
	})

	Context("Given that I have an ASP.NET MVC application (nora) with an application path", func() {
		var app hwcApp

		const contextPath = "/vdir1/vdir2"

		BeforeEach(func() {
			port := newRandomPort()

			env := []string{
				"VCAP_APPLICATION=" + fmt.Sprintf("{ \"application_uris\": [\"localhost:%d%s\"] }", port, contextPath),
				"VCAP_SERVICES={}",
				"PORT=" + strconv.FormatInt(port, 10),
			}
			app = startAppWithEnv("nora", env, false)
			Eventually(app.session).Should(gbytes.Say("Server Started"))
			app.port = port
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app.session).Should(gexec.Exit(0))
		})

		It("runs it on the specified port and path", func() {
			url := fmt.Sprintf("http://localhost:%d%s", app.port, contextPath)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := io.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal(fmt.Sprintf(`"hello i am %s running on http://localhost:%d%s"`,
				"nora", app.port, contextPath)))
		})
	})

	Context("when multiple apps are started by different hwc processes", func() {
		var (
			app1 hwcApp
			app2 hwcApp
		)

		BeforeEach(func() {
			app1 = startApp("nora")
			Eventually(app1.session).Should(gbytes.Say("Server Started"))
			app2 = startApp("nora")
			Eventually(app2.session).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app1)
			Eventually(app1.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app1.session).Should(gexec.Exit(0))

			stopApp(app2)
			Eventually(app2.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app2.session).Should(gexec.Exit(0))
		})

		It("the site name and id should be unique for each app", func() {
			url := fmt.Sprintf("http://localhost:%d/sitename", app1.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))
			body, err := io.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			domainAppId1 := string(body)

			url = fmt.Sprintf("http://localhost:%d/sitename", app2.port)
			res, err = http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))
			body, err = io.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			domainAppId2 := string(body)

			Expect(domainAppId1).NotTo(Equal(domainAppId2))
		})
	})

	Context("The app has an infinite redirect loop", func() {
		var app hwcApp

		BeforeEach(func() {
			app = startApp("stack-overflow")
			Eventually(app.session).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app.session).Should(gexec.Exit(0))
		})

		It("does not get a stack overflow error", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			_, err := http.Get(url)
			Expect(string(app.session.Err.Contents())).NotTo(ContainSubstring("Process is terminated due to StackOverflowException"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stopped after 10 redirects"))
		})
	})

	Context("Given that I have an ASP.NET Classic application", func() {
		var app hwcApp

		BeforeEach(func() {
			app = startApp("asp-classic")
			Eventually(app.session, 10*time.Second).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app.session).Should(gexec.Exit(0))
		})

		It("runs on the specified port", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := io.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal("Hello World!"))
		})

		It("the asp compiled templates directory is valid", func() {
			url := fmt.Sprintf("http://localhost:%d", app.port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			errorMessage := fmt.Sprintf(`Error: The Template Persistent Cache initialization failed for Application Pool 'AppPool%d' because of the following error: Could not create a Disk Cache Sub-directory for the Application Pool. The data may have additional error codes..`, app.port)
			output, err := exec.Command("powershell", "-command", fmt.Sprintf(`(get-eventlog -LogName Application -Message "%s").Message`, errorMessage)).CombinedOutput()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(ContainSubstring("No matches found"))
		})
	})

	Context("my app has troublesome stuff in web.config", func() {
		var app hwcApp

		BeforeEach(func() {
			app = startAppWithEnv("nora", []string{}, true)
			Eventually(app.session).Should(gbytes.Say("Server Started"))
		})

		AfterEach(func() {
			stopApp(app)
			Eventually(app.session).Should(gbytes.Say("Server Shutdown"))
			Eventually(app.session).Should(gexec.Exit(0))
		})

		It("prints out a warning to the user regarding the bad web.config stuff", func() {
			Eventually(app.session.Err).Should(gbytes.Say("Warning: <httpCompression> should not have any attributes but it has nastykey, anotherbadkey"))
			Eventually(app.session.Err).Should(gbytes.Say("Warning: <httpCompression> should not have any child tags other than <staticTypes>" +
				" and <dynamicTypes> but it has <scheme>"))
		})
	})
})

type hwcApp struct {
	session    *gexec.Session
	port       int64
	appDir     string
	profileDir string
}

func stopApp(app hwcApp) {
	d, err := syscall.LoadDLL("kernel32.dll")
	Expect(err).ToNot(HaveOccurred())
	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	Expect(err).ToNot(HaveOccurred())
	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(app.session.Command.Process.Pid))
	Expect(r).ToNot(Equal(0), fmt.Sprintf("GenerateConsoleCtrlEvent: %v\n", err))
	Eventually(app.session).Should(gexec.Exit())

	Eventually(func() error { return os.RemoveAll(app.appDir) }, 10*time.Second, time.Second).Should(Succeed())
	Expect(os.RemoveAll(app.profileDir)).To(Succeed())
}

func startApp(fixtureName string) hwcApp {
	return startAppWithEnv(fixtureName, []string{}, false)
}

func startAppWithEnv(fixtureName string, env []string, badConfigtest bool) hwcApp {
	cmd := exec.Command(hwcBinPath)

	profileDir, err := os.MkdirTemp("", "hwcappprofile")
	Expect(err).ToNot(HaveOccurred())

	appDir, err := os.MkdirTemp("", "hwctestapp")
	Expect(err).ToNot(HaveOccurred())
	wd, err := os.Getwd()
	Expect(err).ToNot(HaveOccurred())
	Expect(copyDirectory(filepath.Join(wd, "fixtures", fixtureName), appDir)).To(Succeed())

	if badConfigtest {
		err := copyFile(filepath.Join(wd, "fixtures", "webconfigs", "Web.config.bad"), filepath.Join(appDir, "Web.config"))
		Expect(err).NotTo(HaveOccurred())
	}

	port := newRandomPort()

	cmd.Env = append([]string{
		"USERPROFILE=" + profileDir,
		"PORT=" + strconv.FormatInt(port, 10),
		"WINDIR=" + os.Getenv("WINDIR"),
		"SYSTEMROOT=" + os.Getenv("SYSTEMROOT"),
		"APP_NAME=" + fixtureName,
	}, env...)
	cmd.Dir = appDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	return hwcApp{
		session:    session,
		port:       port,
		appDir:     appDir,
		profileDir: profileDir,
	}
}

func newRandomPort() int64 {
	const maxPort = 60000
	const minPort = 1025
	return rand.Int63n(maxPort-minPort) + minPort
}

func copyDirectory(srcDir, destDir string) error {
	destExists, _ := fileExists(destDir)
	if !destExists {
		return errors.New("destination dir must exist")
	}

	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(srcDir, f.Name())
		dest := filepath.Join(destDir, f.Name())

		if f.IsDir() {
			err = os.MkdirAll(dest, f.Mode())
			if err != nil {
				return err
			}
			if err := copyDirectory(src, dest); err != nil {
				return err
			}
		} else {
			rc, err := os.Open(src)
			if err != nil {
				return err
			}

			err = writeToFile(rc, dest, f.Mode())
			if err != nil {
				rc.Close()
				return err
			}
			rc.Close()
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func writeToFile(source io.Reader, destFile string, mode os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fh, source)
	if err != nil {
		return err
	}

	return nil
}
func successfulRequest(url string) http.Header {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Err", err.Error())
		os.Exit(1)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Err", err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()
	Expect(res.StatusCode).To(Equal(200))
	io.ReadAll(res.Body)
	return res.Header
}
