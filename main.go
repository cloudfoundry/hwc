// +build windows

package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	_ "runtime/cgo"
	"strconv"
	"syscall"

	cfenv "github.com/cloudfoundry-community/go-cfenv"

	"code.cloudfoundry.org/hwc/contextpath"
	"code.cloudfoundry.org/hwc/hwcconfig"
	"code.cloudfoundry.org/hwc/validator"
	"code.cloudfoundry.org/hwc/webcore"
)

var (
	appRootPath string
	enable32bit bool
)

func init() {
	flag.StringVar(&appRootPath, "appRootPath", ".", "app web root path")
	flag.BoolVar(&enable32bit, "enable32bit", false, "enable 32Bit App On Win64")
}

func main() {
	flag.Parse()

	// spawn 32bit process
	// 2 options to trigger hwc in 32bit mode
	//   1. execute hwc.exe -enable32bit
	//   2. buildpack to choose hwc_x86.exe in final stage (removes the need for this flag)
	if enable32bit && runtime.GOARCH != "386" {
		hwc86Path, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "hwc_x86.exe"))
		checkErr(err)

		hwc86cmd := exec.Command(hwc86Path, "-appRootPath", appRootPath)
		hwc86cmd.Stdout = os.Stdout
		hwc86cmd.Stderr = os.Stderr
		err = hwc86cmd.Run()
		checkErr(err)

		os.Exit(0)
	}

	if os.Getenv("PORT") == "" {
		checkErr(errors.New("Missing PORT environment variable"))
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	checkErr(err)

	rootPath, err := filepath.Abs(appRootPath)
	checkErr(err)

	if os.Getenv("USERPROFILE") == "" {
		checkErr(errors.New("Missing USERPROFILE environment variable"))
	}
	tmpPath, err := filepath.Abs(filepath.Join(os.Getenv("USERPROFILE"), "tmp"))
	checkErr(err)

	err = os.MkdirAll(tmpPath, 0700)
	checkErr(err)

	contextPath := contextpath.Default()
	if cfenv.IsRunningOnCF() {
		appEnv, err := cfenv.Current()
		if err != nil {
			checkErr(fmt.Errorf("Getting current CF environment: %v", err))
		}

		contextPath, err = contextpath.New(appEnv)
		if err != nil {
			checkErr(fmt.Errorf("Getting CF application context path: %v", err))
		}

		fmt.Println(fmt.Sprintf("Context Path %s", contextPath))
	}

	uuid, err := generateUUID()
	if err != nil {
		checkErr(fmt.Errorf("Generating UUID: %v", err))
	}

	err, config := hwcconfig.New(port, rootPath, tmpPath, contextPath, uuid)
	checkErr(err)

	err = validator.ValidateWebConfig(filepath.Join(rootPath, "Web.config"), os.Stderr)
	checkErr(err)

	err, wc := webcore.New()
	checkErr(err)
	defer syscall.FreeLibrary(wc.Handle)

	checkErr(wc.Activate(
		config.ApplicationHostConfigPath,
		config.WebConfigPath,
		config.Instance))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	checkErr(wc.Shutdown(1, config.Instance))
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s\n", err)
		os.Exit(1)
	}
}

func generateUUID() (string, error) {
	const size = 128 / 8
	const format = "%08x-%04x-%04x-%04x-%012x"
	var u [size]byte
	if _, err := io.ReadFull(rand.Reader, u[0:]); err != nil {
		return "", fmt.Errorf("error reading random number generator: %v", err)
	}
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return fmt.Sprintf(format, u[:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
