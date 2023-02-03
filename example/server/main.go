package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vault-thirteen/SFRODB/server"
	"github.com/vault-thirteen/SFRODB/server/settings"
)

func main() {
	cla, err := readCLA()
	mustBeNoError(err)

	var stn *settings.Settings
	stn, err = settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	err = srv.Start()
	mustBeNoError(err)
	fmt.Println("Listening on " + srv.GetDsn())

	appMustBeStopped := make(chan bool, 1)
	waitForQuitSignalFromOS(&appMustBeStopped)
	<-appMustBeStopped

	err = srv.Stop()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func waitForQuitSignalFromOS(quitChan *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-osSignals
		log.Println("quit signal from OS has been received: ", sig)
		*quitChan <- true
	}()
}
