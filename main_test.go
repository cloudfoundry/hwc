// +build windows

package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		cmd.Dir = filepath.Join(wd, "fixtures", env["APP_NAME"])
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
