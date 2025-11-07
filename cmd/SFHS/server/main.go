package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server"
	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server/Settings"
	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println("Using the default configuration file.")
	}

	var stn *Settings.Settings
	stn, err = Settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println("Server is starting ...")
	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		log.Fatal(cerr)
	}
	switch srv.GetWorkMode() {
	case Settings.ServerModeIdHttp:
		fmt.Println("HTTP Server: " + srv.GetListenDsn())
	case Settings.ServerModeIdHttps:
		fmt.Println("HTTPS Server: " + srv.GetListenDsn())
	}
	fmt.Println("DB Client A: " + srv.GetDbDsnA())
	fmt.Println("DB Client B: " + srv.GetDbDsnB())

	serverMustBeStopped := srv.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	log.Println("Stopping the server ...")
	cerr = srv.Stop()
	if cerr != nil {
		log.Println(cerr)
	}
	log.Println("Server was stopped.")
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

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println("quit signal from OS has been received: ", sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
