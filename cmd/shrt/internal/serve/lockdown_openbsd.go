// See LICENSE file for copyright and license details

package serve

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func init() {
	lockdown = func(dbPath string) {
		path, err := filepath.Abs(dbPath)
		if err != nil {
			panic(fmt.Sprint("provided path cannot be made absolute", dbPath))
		}
		err = unix.Unveil(path, "r")
		if err != nil {
			panic(fmt.Sprint("lockdown: %s", err))
		}
		unix.Pledge("stdio rpath dns inet flock", "")
		if err != nil {
			panic(fmt.Sprint("lockdown: %s", err))
		}
	}
}
