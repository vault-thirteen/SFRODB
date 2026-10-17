package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vault-thirteen/auxie/cio"
)

func WaitForQuitSignalFromOS(appMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				cio.LogMessage(Msg_OsQuitSignal, sig)
				*appMustBeStopped <- true
			}
		}
	}()
}
