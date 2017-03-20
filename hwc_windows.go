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

type App struct {
	Instance              string
	Port                  int
	RootPath              string
	TempDirectory         string
	ApplicationHostConfig string
	AspnetConfig          string
	WebConfig             string
}

func init() {
	flag.StringVar(&appRootPath, "appRootPath", ".", "app web root path")
}

func main() {
	flag.Parse()

	wc, err := newWebCore()
	CheckErr(err)
	defer syscall.FreeLibrary(wc.handle)

	if os.Getenv("PORT") == "" {
		Fail(errors.New("Missing PORT environment variable"))
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	CheckErr(err)

	rootPath, err := filepath.Abs(appRootPath)
	CheckErr(err)

	if os.Getenv("USERPROFILE") == "" {
		Fail(errors.New("Missing USERPROFILE environment variable"))
	}
	tmpPath, err := filepath.Abs(filepath.Join(os.Getenv("USERPROFILE"), "tmp"))
	CheckErr(err)

	err = os.MkdirAll(tmpPath, 0700)
	CheckErr(err)

	uuid, err := GenerateUUID()
	if err != nil {
		Fail(fmt.Errorf("Generating UUID: %v", err))
	}

	app := App{
		Instance:      uuid,
		Port:          port,
		RootPath:      rootPath,
		TempDirectory: tmpPath,
	}
	CheckErr(app.configure())

	CheckErr(wc.activate(
		app.ApplicationHostConfig,
		app.WebConfig,
		app.Instance))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	CheckErr(wc.shutdown(1, app.Instance))
}

func newWebCore() (*webCore, error) {
	hwebcore, err := syscall.LoadLibrary(os.ExpandEnv(`${windir}\\\system32\inetsrv\hwebcore.dll`))
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

func (a App) generateApplicationHostConfig() error {
	file, err := os.Create(a.ApplicationHostConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("applicationhost").Parse(ApplicationHostConfig))
	if err := tmpl.Execute(file, a); err != nil {
		return err
	}
	return nil
}

func (a App) generateAspNetConfig() error {
	file, err := os.Create(a.AspnetConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("aspnet").Parse(AspnetConfig))
	if err := tmpl.Execute(file, a); err != nil {
		return err
	}
	return nil
}

func (a App) generateWebConfig() error {
	file, err := os.Create(a.WebConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpl = template.Must(template.New("webconfig").Parse(WebConfig))
	if err := tmpl.Execute(file, a); err != nil {
		return err
	}
	return nil
}

func (a *App) configure() error {
	dest := filepath.Join(a.TempDirectory, "config")
	err := os.MkdirAll(dest, 0700)
	if err != nil {
		return err
	}

	a.ApplicationHostConfig = filepath.Join(dest, "ApplicationHost.config")
	a.AspnetConfig = filepath.Join(dest, "Aspnet.config")
	a.WebConfig = filepath.Join(dest, "Web.config")

	err = a.generateApplicationHostConfig()
	if err != nil {
		return err
	}

	err = a.generateAspNetConfig()
	if err != nil {
		return err
	}

	err = a.generateWebConfig()
	if err != nil {
		return err
	}

	return nil
}

func CheckErr(err error) {
	if err != nil {
		Fail(err)
	}
}

func Fail(err error) {
	fmt.Fprintf(os.Stderr, "\n%s\n", err)
	os.Exit(1)
}

func GenerateUUID() (string, error) {
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
