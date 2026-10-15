package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"

	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server"
	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server/Settings"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	Msg_ServerIsStarting = "Server is starting ..."
	Msg_StoppingServer   = "Stopping the server ..."
	Msg_ServerStopped    = "Server was stopped."
	Msg_HttpServer       = "HTTP Server: "
	Msg_HttpsServer      = "HTTPS Server: "
	Msg_DbClientA        = "DB Client A: "
	Msg_DbClientB        = "DB Client B: "
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

	var stn *Settings.Settings
	stn, err = Settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println(Msg_ServerIsStarting)
	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		log.Fatal(cerr)
	}
	switch srv.GetWorkMode() {
	case Settings.ServerModeIdHttp:
		fmt.Println(Msg_HttpServer + srv.GetListenDsn())
	case Settings.ServerModeIdHttps:
		fmt.Println(Msg_HttpsServer + srv.GetListenDsn())
	}
	fmt.Println(Msg_DbClientA + srv.GetDbDsnA())
	fmt.Println(Msg_DbClientB + srv.GetDbDsnB())

	serverMustBeStopped := srv.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	log.Println(Msg_StoppingServer)
	cerr = srv.Stop()
	if cerr != nil {
		log.Println(cerr)
	}
	log.Println(Msg_ServerStopped)
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

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println(common.Msg_OsQuitSignal, sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
