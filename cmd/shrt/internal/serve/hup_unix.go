// See LICENSE file for copyright and license details

package serve

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"djmo.ch/go-shrt"
)

func init() {
	hangup = func(h *shrt.ShrtHandler) {
		mux := h.GetMutex()
		hup := make(chan os.Signal, 1)
		signal.Notify(hup, syscall.SIGHUP)
		for {
			<-hup
			tmpShrt, err := shrt.ReadShrtFile(h.Config.DbPath)
			if err != nil {
				log.Println("db error:", err)
			} else {
				mux.Lock()
				h.ShrtFile = tmpShrt
				mux.Unlock()
			}
		}
	}
}
