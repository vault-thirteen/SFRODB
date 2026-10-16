package main

import (
	"math"
	"os"
	"strconv"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/random"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommandLineArguments/client"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

func main() {
	showIntro()

	clargs, err := clac.NewCommandLineArgumentsFromOs()
	mustBeNoError(err)

	var stn *cs.ClientSettings
	stn, err = cs.New(clargs.Host, clargs.MainPort, clargs.AuxPort, 0)
	mustBeNoError(err)
	common.LogMessage(common.Msg_Settings, stn)

	var clientId uint
	clientId, err = random.Uint(1, math.MaxUint32)
	mustBeNoError(err)

	var cli *client.Client
	cli, err = client.New(stn, strconv.FormatUint(uint64(clientId), 10))
	mustBeNoError(err)

	cerr := cli.Start()
	if cerr != nil {
		errors.LogError(cerr)
		os.Exit(common.ExitCode_StartError)
		return
	}
	common.PrintMessage(common.Msg_ConnectedTo + cli.GetMainDsn())
	common.PrintMessage(common.Msg_ConnectedTo + cli.GetAuxDsn())

	appMustBeStopped := make(chan bool, 1)
	go makeSomeActions(cli, &appMustBeStopped)
	common.WaitForQuitSignalFromOS(&appMustBeStopped)
	<-appMustBeStopped

	cerr = cli.Stop()
	if cerr != nil {
		errors.LogError(cerr)
	}
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
	versioneer.ShowIntroText(common.App_TestClient)
	versioneer.ShowComponentsInfoText()
	common.PrintMessage()
}
