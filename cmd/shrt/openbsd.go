// See LICENSE file for copyright and license details
// +build openbsd

package main

// #include <stdlib.h>
// #include <unistd.h>
import "C"

import (
	"fmt"
	"path/filepath"
	"unsafe"
)

func init() {
	osInit = func(dbPath string) error {
		path, err := filepath.Abs(dbPath)
		if err != nil {
			return fmt.Errorf("osInit: provided path cannot be made absolute")
		}
		err = unveil(path, "r")
		if err != nil {
			return fmt.Errorf("osInit: %s", err.Error())
		}
		pledge("stdio rpath inet flock", "")
		if err != nil {
			return fmt.Errorf("osInit: %s", err.Error())
		}
		return nil
	}
}

func pledge(promises, execpromises string) error {
	cPromises := C.CString(promises)
	defer C.free(unsafe.Pointer(cPromises))
	cExecPromises := C.CString(execpromises)
	defer C.free(unsafe.Pointer(cExecPromises))

	if eVal, err := C.pledge(cPromises, cExecPromises); eVal != 0 {
		return fmt.Errorf("pledge: %d", err)
	}

	return nil
}

func unveil(path, permissions string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cPermissions := C.CString(permissions)
	defer C.free(unsafe.Pointer(cPermissions))

	if eVal, err := C.unveil(cPath, cPermissions); eVal != 0 {
		return fmt.Errorf("unveil: %d", err)
	}

	return nil
}
