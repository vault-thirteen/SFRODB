package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vault-thirteen/SFRODB/client"
	"github.com/vault-thirteen/SFRODB/client/settings"
)

func main() {
	cla, err := readCLA()
	mustBeNoError(err)

	var stn *settings.Settings
	stn, err = settings.NewSettings(cla.Host, cla.MainPort, cla.AuxPort, 0)
	mustBeNoError(err)
	log.Println("Settings:", stn)

	var cli *client.Client
	cli, err = client.NewClient(stn, "1")
	mustBeNoError(err)

	cerr := cli.Start()
	if cerr != nil {
		log.Fatal(cerr)
	}
	fmt.Println("Connected to " + cli.GetMainDsn())
	fmt.Println("Connected to " + cli.GetAuxDsn())

	appMustBeStopped := make(chan bool, 1)
	go makeSomeActions(cli, &appMustBeStopped)

	waitForQuitSignalFromOS(&appMustBeStopped)
	<-appMustBeStopped

	cerr = cli.Stop()
	if cerr != nil {
		log.Println(cerr)
	}
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
