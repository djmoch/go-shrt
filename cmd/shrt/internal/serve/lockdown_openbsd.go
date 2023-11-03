// See LICENSE file for copyright and license details

package serve

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func init() {
	lockdown = func(dbPath string) {
		path := "/" + dbPath
		err := unix.Unveil(path, "r")
		if err != nil {
			panic(fmt.Sprint("lockdown: ", err))
		}
		err = unix.Pledge("stdio rpath dns inet flock", "")
		if err != nil {
			panic(fmt.Sprint("lockdown: ", err))
		}
	}
}
