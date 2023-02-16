package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cache "github.com/vault-thirteen/Cache"
	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/server"
	"github.com/vault-thirteen/SFRODB/server/settings"
)

func main() {
	common.ShowIntroText(common.ProductServer)
	common.ShowComponentInfoText(common.ComponentCache, cache.LibVersion)

	cla, err := readCLA()
	mustBeNoError(err)

	var stn *settings.Settings
	stn, err = settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		log.Fatal(cerr)
	}
	fmt.Println("Main Listener: " + srv.GetMainDsn())
	fmt.Println("Auxiliary Listener: " + srv.GetAuxDsn())

	appMustBeStopped := make(chan bool, 1)
	waitForQuitSignalFromOS(&appMustBeStopped)
	<-appMustBeStopped

	cerr = srv.Stop()
	if cerr != nil {
		log.Println(cerr)
	}
	time.Sleep(time.Second)
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
