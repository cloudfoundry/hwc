//go:build windows
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

		appHostConfigPathPtr, err := syscall.UTF16PtrFromString(appHostConfigPath)
		if err != nil {
			return err
		}
		rootWebConfigPathPtr, err := syscall.UTF16PtrFromString(rootWebConfigPath)
		if err != nil {
			return err
		}
		instanceNamePtr, err := syscall.UTF16PtrFromString(instanceName)
		if err != nil {
			return err
		}

		r1, _, exitCode := syscall.SyscallN(uintptr(webCoreActivate),
			uintptr(unsafe.Pointer(appHostConfigPathPtr)),
			uintptr(unsafe.Pointer(rootWebConfigPathPtr)),
			uintptr(unsafe.Pointer(instanceNamePtr)))
		if exitCode != 0 {
			return fmt.Errorf("WebCoreActivate returned exit code: %d", exitCode)
		}
		if r1 != 0 {
			return fmt.Errorf("HWC Failed to start: return code: 0x%02x", r1)
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

		_, _, exitCode := syscall.SyscallN(uintptr(webCoreShutdown),
			uintptr(unsafe.Pointer(&immediate)), 0, 0)
		if exitCode != 0 {
			return fmt.Errorf("WebCoreShutdown returned exit code: %d", exitCode)
		}
		fmt.Printf("Server Shutdown for %+v\n", instanceName)
	}

	return nil
}
