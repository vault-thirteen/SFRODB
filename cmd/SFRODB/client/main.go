package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/random"
)

func main() {
	showIntro()

	cla, err := readCLA()
	mustBeNoError(err)

	var stn *cs.ClientSettings
	stn, err = cs.New(cla.Host, cla.MainPort, cla.AuxPort, 0)
	mustBeNoError(err)
	log.Println("Settings:", stn)

	var clientId uint
	clientId, err = random.Uint(1, math.MaxUint32)
	mustBeNoError(err)

	var cli *client.Client
	cli, err = client.New(stn, strconv.FormatUint(uint64(clientId), 10))
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

func showIntro() {
	versioneer, err := ver.New()
	mustBeNoError(err)
	versioneer.ShowIntroText("Test Client")
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
