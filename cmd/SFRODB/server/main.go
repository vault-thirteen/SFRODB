package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Server"
	ss "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ServerSettings"
	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println("Using the default configuration file.")
	}

	var stn *ss.ServerSettings
	stn, err = ss.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	var srv *server.Server
	srv, err = server.New(stn)
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

func showIntro() {
	versioneer, err := ver.New()
	mustBeNoError(err)
	versioneer.ShowIntroText("Server")
	versioneer.ShowComponentsInfoText()
	fmt.Println()
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
