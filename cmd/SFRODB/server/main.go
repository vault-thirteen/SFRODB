package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Server"
	ss "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ServerSettings"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	Msg_MainListener      = "Main Listener: "
	Msg_AuxiliaryListener = "Auxiliary Listener: "
)

const (
	App = "Server"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println(common.Msg_UsingDefaultConfigurationFile)
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
	fmt.Println(Msg_MainListener + srv.GetMainDsn())
	fmt.Println(Msg_AuxiliaryListener + srv.GetAuxDsn())

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
	versioneer, err := ver.New(false)
	mustBeNoError(err)
	versioneer.ShowIntroText(App)
	versioneer.ShowComponentsInfoText()
	fmt.Println()
}

func waitForQuitSignalFromOS(quitChan *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-osSignals
		log.Println(common.Msg_OsQuitSignal, sig)
		*quitChan <- true
	}()
}
