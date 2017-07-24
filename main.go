// +build windows

package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	_ "runtime/cgo"
	"strconv"
	"syscall"

	cfenv "github.com/cloudfoundry-community/go-cfenv"

	"code.cloudfoundry.org/hwc/contextpath"
	"code.cloudfoundry.org/hwc/hwcconfig"
	"code.cloudfoundry.org/hwc/webcore"
)

var appRootPath string

func init() {
	flag.StringVar(&appRootPath, "appRootPath", ".", "app web root path")
}

func main() {
	flag.Parse()

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
