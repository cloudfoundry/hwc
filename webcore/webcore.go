// +build windows

package webcore

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type WebCore struct {
	activated bool
	Handle    syscall.Handle
}

func New() (error, *WebCore) {
	hwebcore, err := syscall.LoadLibrary(os.ExpandEnv(`${windir}\system32\inetsrv\hwebcore.dll`))
	if err != nil {
		return err, nil
	}

	return nil, &WebCore{
		activated: false,
		Handle:    hwebcore,
	}
}

func (w *WebCore) Activate(appHostConfigPath, rootWebConfigPath, instanceName string) error {
	if !w.activated {
		webCoreActivate, err := syscall.GetProcAddress(w.Handle, "WebCoreActivate")
		if err != nil {
			return err
		}

		var nargs uintptr = 3
		_, _, exitCode := syscall.Syscall(uintptr(webCoreActivate),
			nargs,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(appHostConfigPath))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(rootWebConfigPath))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(instanceName))))
		if exitCode != 0 {
			return fmt.Errorf("WebCoreActivate returned exit code: %d", exitCode)
		}

		fmt.Printf("Server Started for %+v\n", instanceName)
		w.activated = true
	}

	return nil
}

func (w *WebCore) Shutdown(immediate int, instanceName string) error {
	if w.activated {
		webCoreShutdown, err := syscall.GetProcAddress(w.Handle, "WebCoreShutdown")
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
