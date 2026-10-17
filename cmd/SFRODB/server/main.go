package main

import (
	"os"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommandLineArguments/server"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Server"
	ss "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ServerSettings"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

func main() {
	showIntro()

	cla, err := clas.NewCommandLineArgumentsFromOs()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		cio.LogMessage(common.Msg_UsingDefaultConfigurationFile)
	}

	var stn *ss.ServerSettings
	stn, err = ss.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	var srv *server.Server
	srv, err = server.New(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		errors.LogError(cerr)
		os.Exit(common.ExitCode_StartError)
		return
	}
	cio.PrintMessage(common.Msg_ListenerMain + srv.GetMainDsn())
	cio.PrintMessage(common.Msg_ListenerAuxiliary + srv.GetAuxDsn())

	appMustBeStopped := make(chan bool, 1)
	common.WaitForQuitSignalFromOS(&appMustBeStopped)
	<-appMustBeStopped

	cerr = srv.Stop()
	if cerr != nil {
		errors.LogError(cerr)
	}
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
