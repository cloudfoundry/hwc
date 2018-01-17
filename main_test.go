// +build windows

package main_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("HWC", func() {
	var (
		err        error
		binaryPath string
		tmpDir     string
		env        map[string]string
	)

	BeforeEach(func() {
		binaryPath, err = BuildWithEnvironment("code.cloudfoundry.org/hwc", []string{"CGO_ENABLED=1", "GO_EXTLINK_ENABLED=1"})
		Expect(err).ToNot(HaveOccurred())
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
		env = map[string]string{
			"USERPROFILE": tmpDir,
			"PORT":        "43311",
			"WINDIR":      os.Getenv("WINDIR"),
			"SYSTEMROOT":  os.Getenv("SYSTEMROOT"),
			"APP_NAME":    "nora",
		}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
		CleanupBuildArtifacts()
	})

	sendCtrlBreak := func(s *Session) {
		d, err := syscall.LoadDLL("kernel32.dll")
		Expect(err).ToNot(HaveOccurred())
		p, err := d.FindProc("GenerateConsoleCtrlEvent")
		Expect(err).ToNot(HaveOccurred())
		r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(s.Command.Process.Pid))
		Expect(r).ToNot(Equal(0), fmt.Sprintf("GenerateConsoleCtrlEvent: %v\n", err))
	}

	startApp := func(env map[string]string) (*Session, error) {
		cmd := exec.Command(binaryPath)
		vals := []string{}
		for k, v := range env {
			vals = append(vals, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = vals
		cmd.Dir = env["APP_DIR"]
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}

		return Start(cmd, GinkgoWriter, GinkgoWriter)
	}

	Context("when the app PORT is not set", func() {
		JustBeforeEach(func() {
			env["PORT"] = ""
		})

		It("errors", func() {
			session, err := startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Exit(1))
			Eventually(session.Err).Should(Say("Missing PORT environment variable"))
		})
	})

	Context("when the app USERPROFILE is not set", func() {
		JustBeforeEach(func() {
			env["USERPROFILE"] = ""
		})

		It("errors", func() {
			session, err := startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Exit(1))
			Eventually(session.Err).Should(Say("Missing USERPROFILE environment variable"))
		})
	})

	Context("Given that I am missing a required .dll", func() {
		JustBeforeEach(func() {
			env["WINDIR"] = tmpDir
		})

		It("errors", func() {
			session, err := startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Exit(1))
			Eventually(session.Err).Should(Say("Missing required DLLs:"))
		})
	})

	Context("Given that I have an ASP.NET MVC application (nora)", func() {
		var (
			session *Session
			err     error
		)

		BeforeEach(func() {
			session, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Say("Server Started"))
		})

		AfterEach(func() {
			sendCtrlBreak(session)
			Eventually(session, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(session).Should(Exit(0))
		})

		It("runs it on the specified port", func() {
			url := fmt.Sprintf("http://localhost:%s", env["PORT"])
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal(fmt.Sprintf(`"hello i am %s running on http://localhost:%s/"`, env["APP_NAME"], env["PORT"])))
		})

		It("correctly utilizes the USERPROFILE temp directory", func() {
			url := fmt.Sprintf("http://localhost:%s", env["PORT"])
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			Expect(filepath.Join(tmpDir, "tmp", "root")).To(BeADirectory())
			Expect(err).ToNot(HaveOccurred())

			By("placing config files in the temp directory", func() {
				Expect(filepath.Join(tmpDir, "tmp", "config", "Web.config")).To(BeAnExistingFile())
				Expect(filepath.Join(tmpDir, "tmp", "config", "ApplicationHost.config")).To(BeAnExistingFile())
				Expect(filepath.Join(tmpDir, "tmp", "config", "Aspnet.config")).To(BeAnExistingFile())
			})
		})

		It("does not add unexpected custom headers", func() {
			url := fmt.Sprintf("http://localhost:%s", env["PORT"])
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			var customHeaders []string
			for h, _ := range res.Header {
				if strings.HasPrefix(h, "X-") {
					customHeaders = append(customHeaders, strings.ToLower(h))
				}
			}
			Expect(len(customHeaders)).To(Equal(1))
			Expect(customHeaders[0]).To(Equal("x-aspnet-version"))
		})
	})

	Context("Given that I have an ASP.NET MVC application (nora) with an application path", func() {
		var (
			session *Session
			err     error
		)

		const contextPath = "/vdir1/vdir2"

		BeforeEach(func() {
			env["VCAP_APPLICATION"] = fmt.Sprintf("{ \"application_uris\": [\"localhost:%s%s\"] }", env["PORT"], contextPath)
			env["VCAP_SERVICES"] = "{}"
			session, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Say("Server Started"))
		})

		AfterEach(func() {
			sendCtrlBreak(session)
			Eventually(session, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(session).Should(Exit(0))
		})

		It("runs it on the specified port and path", func() {
			url := fmt.Sprintf("http://localhost:%s%s", env["PORT"], contextPath)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal(fmt.Sprintf(`"hello i am %s running on http://localhost:%s%s"`,
				env["APP_NAME"], env["PORT"], contextPath)))
		})
	})

	Context("when multiple apps are started by different hwc processes", func() {
		var (
			app1        *Session
			app1Port    string
			app1Profile string
			app1Dir     string

			app2        *Session
			app2Port    string
			app2Profile string
			app2Dir     string
		)

		BeforeEach(func() {
			var err error
			wd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			app1Port = newRandomPortStr()
			env["PORT"] = app1Port

			app1Profile, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())
			env["USERPROFILE"] = app1Profile

			app1Dir, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())
			env["APP_DIR"] = app1Dir
			Expect(copyDirectory(filepath.Join(wd, "fixtures", "nora"), app1Dir)).To(Succeed())

			app1, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(app1, 10*time.Second).Should(Say("Server Started"))

			app2Port = newRandomPortStr()
			env["PORT"] = app2Port

			app2Profile, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())
			env["USERPROFILE"] = app2Profile

			app2Dir, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())
			env["APP_DIR"] = app2Dir
			Expect(copyDirectory(filepath.Join(wd, "fixtures", "nora"), app2Dir)).To(Succeed())

			app2, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(app2, 10*time.Second).Should(Say("Server Started"))
		})

		AfterEach(func() {
			sendCtrlBreak(app1)
			Eventually(app1, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(app1).Should(Exit(0))
			Expect(os.RemoveAll(app1Profile)).To(Succeed())
			Expect(os.RemoveAll(app1Dir)).To(Succeed())

			sendCtrlBreak(app2)
			Eventually(app2, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(app2).Should(Exit(0))
			Expect(os.RemoveAll(app2Profile)).To(Succeed())
			Expect(os.RemoveAll(app2Dir)).To(Succeed())
		})

		FIt("the site name and id should be unique for each app", func() {
			url := fmt.Sprintf("http://localhost:%s/sitename", app1Port)
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200), string(body))

			// body, err := ioutil.ReadAll(res.Body)
			// Expect(err).ToNot(HaveOccurred())
			app1DomainAppId := string(body)

			url = fmt.Sprintf("http://localhost:%s/sitename", app2Port)
			res, err = http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err = ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			app2DomainAppId := string(body)

			Expect(app1DomainAppId).NotTo(Equal(app2DomainAppId))
		})
	})

	Context("The app has an infinite redirect loop", func() {
		var (
			session *Session
			err     error
		)

		BeforeEach(func() {
			env["APP_NAME"] = "stack-overflow"
			session, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Say("Server Started"))
		})

		AfterEach(func() {
			sendCtrlBreak(session)
			Eventually(session, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(session).Should(Exit(0))
		})

		It("does not get a stack overflow error", func() {
			url := fmt.Sprintf("http://localhost:%s", env["PORT"])
			_, err := http.Get(url)
			Expect(string(session.Err.Contents())).NotTo(ContainSubstring("Process is terminated due to StackOverflowException"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stopped after 10 redirects"))
		})
	})

	Context("Given that I have an ASP.NET Classic application", func() {
		var (
			session *Session
			err     error
		)

		BeforeEach(func() {
			env["APP_NAME"] = "asp-classic"
			session, err = startApp(env)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 10*time.Second).Should(Say("Server Started"))
		})

		AfterEach(func() {
			sendCtrlBreak(session)
			Eventually(session, 10*time.Second).Should(Say("Server Shutdown"))
			Eventually(session).Should(Exit(0))
		})

		It("runs on the specified port", func() {
			url := fmt.Sprintf("http://localhost:%s", env["PORT"])
			res, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(res.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(Equal("Hello World!"))
		})
	})
})

func newRandomPortStr() string {
	const maxPort = 60000
	const minPort = 1025
	port := rand.Int63n(maxPort-minPort) + minPort
	return strconv.FormatInt(port, 10)
}

func copyDirectory(srcDir, destDir string) error {
	destExists, _ := fileExists(destDir)
	if !destExists {
		return errors.New("destination dir must exist")
	}

	files, err := ioutil.ReadDir(srcDir)
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
