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
	"strconv"
	"syscall"
	"text/template"
	"unsafe"
)

var appRootPath string

type webCore struct {
	activated bool
	handle    syscall.Handle
}

type HwcConfig struct {
	Instance      string
	Port          int
	RootPath      string
	TempDirectory string

	AspnetConfigPath string
	WebConfigPath    string

	ApplicationHostConfigPath string
}

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

	uuid, err := generateUUID()
	if err != nil {
		checkErr(fmt.Errorf("Generating UUID: %v", err))
	}

	err, config := NewHwcConfig(uuid, rootPath, tmpPath, port)
	checkErr(err)

	wc, err := newWebCore()
	checkErr(err)
	defer syscall.FreeLibrary(wc.handle)

	checkErr(wc.activate(
		config.ApplicationHostConfigPath,
		config.WebConfigPath,
		config.Instance))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	checkErr(wc.shutdown(1, config.Instance))
}

func NewHwcConfig(uuid, rootPath, tmpPath string, port int) (error, *HwcConfig) {
	config := &HwcConfig{
		Instance:      uuid,
		Port:          port,
		RootPath:      rootPath,
		TempDirectory: tmpPath,
	}
	dest := filepath.Join(config.TempDirectory, "config")
	err := os.MkdirAll(dest, 0700)
	if err != nil {
		return err, nil
	}

	config.ApplicationHostConfigPath = filepath.Join(dest, "ApplicationHost.config")
	config.AspnetConfigPath = filepath.Join(dest, "Aspnet.config")
	config.WebConfigPath = filepath.Join(dest, "Web.config")

	err, applicationHostConfig := NewApplicationHostConfig()
	if err != nil {
		return err, nil
	}
	err = applicationHostConfig.generate(*config)
	if err != nil {
		return err, nil
	}

	err = config.generateAspNetConfig()
	if err != nil {
		return err, nil
	}

	err = config.generateWebConfig()
	if err != nil {
		return err, nil
	}

	return nil, config
}

func newWebCore() (*webCore, error) {
	hwebcore, err := syscall.LoadLibrary(os.ExpandEnv(`${windir}\system32\inetsrv\hwebcore.dll`))
	if err != nil {
		return nil, err
	}

	return &webCore{
		activated: false,
		handle:    hwebcore,
	}, nil
}

func (w *webCore) activate(appHostConfig, rootWebConfig, instanceName string) error {
	if !w.activated {
		webCoreActivate, err := syscall.GetProcAddress(w.handle, "WebCoreActivate")
		if err != nil {
			return err
		}

		var nargs uintptr = 3
		_, _, exitCode := syscall.Syscall(uintptr(webCoreActivate),
			nargs,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(appHostConfig))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(rootWebConfig))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(instanceName))))
		if exitCode != 0 {
			return fmt.Errorf("WebCoreActivate returned exit code: %d", exitCode)
		}

		fmt.Printf("Server Started for %+v\n", instanceName)
		w.activated = true
	}

	return nil
}

func (w *webCore) shutdown(immediate int, instanceName string) error {
	if w.activated {
		webCoreShutdown, err := syscall.GetProcAddress(w.handle, "WebCoreShutdown")
		if err != nil {
			return err
		}

		var nargs uintptr = 1
		_, _, exitCode := syscall.Syscall(uintptr(webCoreShutdown),
			nargs, uintptr(unsafe.Pointer(&immediate)), 0, 0)
		if exitCode != 0 {
			return fmt.Errorf("WebCoreShutdown returned exit code: %d", exitCode)
		}
		fmt.Printf("Server Shutdown for %+v\n", instanceName)
	}

	return nil
}

func (hc *HwcConfig) generateAspNetConfig() error {
	file, err := os.Create(hc.AspnetConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("aspnet").Parse(AspnetConfig))
	if err := tmpl.Execute(file, hc); err != nil {
		return err
	}
	return nil
}

func (hc *HwcConfig) generateWebConfig() error {
	file, err := os.Create(hc.WebConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("webconfig").Parse(WebConfig))
	if err := tmpl.Execute(file, hc); err != nil {
		return err
	}
	return nil
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
