// See LICENSE file for copyright and license details

package serve

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"djmo.ch/go-shrt"
)

func init() {
	hangup = func(h *shrt.ShrtHandler) {
		hup := make(chan os.Signal, 1)
		signal.Notify(hup, syscall.SIGHUP)
		for {
			<-hup
			f, err := h.FS.Open(h.Config.DbPath)
			if err != nil {
				panic(fmt.Sprint("could not open", h.Config.DbPath))
			}
			err = h.ShrtFile.ReadShrtFile(f)
			if err != nil {
				panic(fmt.Sprint("db error:", err))
			}
		}
	}
}
