package main

import (
	"os"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/auxie/cio"

	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server"
	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server/CommandLineArguments"
	"github.com/vault-thirteen/SFRODB/pkg/SFHS/server/Settings"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

func main() {
	showIntro()

	clargs, err := cla.NewCommandLineArgumentsFromOs()
	mustBeNoError(err)
	if clargs.IsDefaultFile() {
		cio.LogMessage(common.Msg_UsingDefaultConfigurationFile)
	}

	var stn *Settings.Settings
	stn, err = Settings.NewSettingsFromFile(clargs.ConfigurationFilePath)
	mustBeNoError(err)

	cio.LogMessage(common.Msg_ServerStarting)
	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		errors.LogError(cerr)
		os.Exit(common.ExitCode_StartError)
		return
	}
	switch srv.GetWorkMode() {
	case Settings.ServerModeId_Http:
		cio.PrintMessage(common.Msg_HttpServer + srv.GetListenDsn())
	case Settings.ServerModeId_Https:
		cio.PrintMessage(common.Msg_HttpsServer + srv.GetListenDsn())
	}
	cio.PrintMessage(common.Msg_DbClientA + srv.GetDbDsnA())
	cio.PrintMessage(common.Msg_DbClientB + srv.GetDbDsnB())

	serverMustBeStopped := srv.GetStopChannel()
	common.WaitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	cio.LogMessage(common.Msg_ServerStopping)
	cerr = srv.Stop()
	if cerr != nil {
		errors.LogError(cerr)
	}
	cio.LogMessage(common.Msg_ServerStopped)
	time.Sleep(time.Second)
}

func mustBeNoError(err error) {
	if err != nil {
		errors.LogError(err)
		os.Exit(common.ExitCode_Init)
		return
	}
}

func showIntro() {
	versioneer, err := ver.New(false)
	mustBeNoError(err)
	versioneer.ShowIntroText(common.App_Server)
	versioneer.ShowComponentsInfoText()
	cio.PrintMessage()
}
