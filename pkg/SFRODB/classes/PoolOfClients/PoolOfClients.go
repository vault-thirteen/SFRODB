package poc

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vault-thirteen/auxie/cio"
	ae "github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
)

const (
	ClientRestarterSuccessPauseSec = 5
	ClientRestarterFailurePauseSec = 15
)

const (
	Err_NoIdleClientIsAvailable = "no idle client is available"

	Errf_DuplicateClientId    = "duplicate client id: %v"
	Errf_ClientIsNotBeingUsed = "client is not being used: %v"
)

const (
	Msg_BrokenClientReconnected     = "A broken client was successfully reconnected."
	Msg_ClientRestarterStopped      = "Client restarter has stopped."
	Msg_ClientPoolStopped           = "Client pool has been stopped."
	Msg_ReconnectingClient          = "Re-connecting the client ..."
	Msg_StoppingAllClients          = "Stopping all the clients ..."
	Msg_WaitingForSubroutinesToStop = "Waiting for subroutines to stop ..."

	Msgf_ClientsPoolStart = "A pool of %d clients has been started."
	Msgf_ClientStopped    = "Client [%s] was stopped."
)

type PoolOfClients struct {
	size int

	// Idle clients.
	// Clients with established connection, ready to be used.
	idleClients chan *client.Client

	// Clients which are actively used at the moment.
	usedClients map[string]*client.Client

	// Clients which had errors are put into the 'brokenClients' channel.
	// Client restarter periodically tries to restart them.
	brokenClients chan *client.Client

	// This channel may be used only by the 'Stop' method.
	// No one else is allowed to use this channel.
	stoppedClients chan *client.Client

	// Control structures.
	clientTransfers *sync.Mutex
	subRoutines     *sync.WaitGroup
	mustStop        *atomic.Bool
}

func New(size int, stn *cs.ClientSettings) (pool *PoolOfClients, err error) {
	pool = &PoolOfClients{
		size:            size,
		idleClients:     make(chan *client.Client, size),
		usedClients:     make(map[string]*client.Client),
		brokenClients:   make(chan *client.Client, size),
		stoppedClients:  make(chan *client.Client, size),
		clientTransfers: new(sync.Mutex),
		subRoutines:     new(sync.WaitGroup),
		mustStop:        new(atomic.Bool),
	}
	pool.mustStop.Store(false)

	var cli *client.Client
	for i := 1; i <= size; i++ {
		cli, err = client.New(stn, strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		pool.stoppedClients <- cli
	}

	return pool, nil
}

func (cp *PoolOfClients) Start() (cerr *ce.CommonError) {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	// Take all the clients.
	clientsToStart := make([]*client.Client, 0, cp.size)
	for i := 1; i <= cp.size; i++ {
		cli := <-cp.stoppedClients
		clientsToStart = append(clientsToStart, cli)
	}

	// If at least one client fails to start, we must:
	//	- stop all the clients;
	//	- put all the clients into the channel as before the start.
	defer func() {
		if cerr != nil {
			for _, cli := range clientsToStart {
				_ = cli.Stop()
				cp.stoppedClients <- cli
			}
		}
	}()

	// Start all the clients.
	for _, cli := range clientsToStart {
		cerr = cli.Start()
		if cerr != nil {
			return cerr
		}
	}

	// Save all the started clients.
	for _, cli := range clientsToStart {
		cp.idleClients <- cli
	}

	cp.subRoutines.Add(1)
	go cp.clientRestarter()

	cio.LogMessage(cio.FormatMessage(Msgf_ClientsPoolStart, cp.size))

	return nil
}

// clientRestarter periodically tries to restart broken clients.
func (cp *PoolOfClients) clientRestarter() {
	defer cp.subRoutines.Done()

	var ok bool
	for {
		if cp.mustStop.Load() {
			break
		}

		ok = cp.clientRestarterCore()
		if ok {
			// Wait after success.
			for i := 1; i <= ClientRestarterSuccessPauseSec; i++ {
				if cp.mustStop.Load() {
					break
				}
				time.Sleep(time.Second)
			}
			continue
		}

		// Wait after failure.
		for i := 1; i <= ClientRestarterFailurePauseSec; i++ {
			if cp.mustStop.Load() {
				break
			}
			time.Sleep(time.Second)
		}
		continue
	}

	cio.LogMessage(Msg_ClientRestarterStopped)
}

// clientRestarterCore is a critical part of the client restarter which
// requires mutex. It tries to get a single broken client and restart it. A
// successfully restarted broken client is moved to the 'idleClients' channel.
// Returns 'true' on success.
func (cp *PoolOfClients) clientRestarterCore() (success bool) {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	if len(cp.brokenClients) == 0 {
		return true
	}

	cli := <-cp.brokenClients
	cio.LogMessage(Msg_ReconnectingClient)
	_ = cli.Stop()
	cerr := cli.Start()
	if cerr != nil {
		cp.brokenClients <- cli
		return false
	}

	cio.LogMessage(Msg_BrokenClientReconnected)
	cp.idleClients <- cli
	return true
}

func (cp *PoolOfClients) Stop() {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	cp.mustStop.Store(true)

	// Drain all the channels except the 'stoppedClients' channel and stop all
	// the clients.
	cio.LogMessage(Msg_StoppingAllClients)
	var cli *client.Client
	for len(cp.idleClients) > 0 {
		cli = <-cp.idleClients
		_ = cli.Stop()
		cio.LogMessage(cio.FormatMessage(Msgf_ClientStopped, cli.GetId()))
		cp.stoppedClients <- cli
	}
	for i := range cp.usedClients {
		cli = cp.usedClients[i]
		_ = cli.Stop()
		cio.LogMessage(cio.FormatMessage(Msgf_ClientStopped, cli.GetId()))
		cp.stoppedClients <- cli
		delete(cp.usedClients, cli.GetId())
	}
	for len(cp.brokenClients) > 0 {
		cli = <-cp.brokenClients
		_ = cli.Stop()
		cio.LogMessage(cio.FormatMessage(Msgf_ClientStopped, cli.GetId()))
		cp.stoppedClients <- cli
	}

	cio.LogMessage(Msg_WaitingForSubroutinesToStop)
	cp.subRoutines.Wait()

	cio.LogMessage(Msg_ClientPoolStopped)
}

// GiveIdleClient provides an idle client. If no clients are idle, an error is
// returned.
func (cp *PoolOfClients) GiveIdleClient() (cli *client.Client, err error) {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	if len(cp.idleClients) == 0 {
		return nil, ae.ErrorFromText(Err_NoIdleClientIsAvailable)
	}

	cli = <-cp.idleClients
	_, idExists := cp.usedClients[cli.GetId()]
	if idExists {
		// Clients may not have same IDs !
		cp.idleClients <- cli
		return nil, ae.ErrorFromFormattedText(Errf_DuplicateClientId, cli.GetId())
	}

	cp.usedClients[cli.GetId()] = cli
	return cli, nil
}

// TakeIdleClient receives an idle client.
// While current version of the Go programming language does not allow to use
// methods applicable to older versions of the language, now we can not know
// whether the client is functioning or broken. In previous versions of Go
// it was possible to know the state of the network connection by reading it
// into an empty slice. By the way, TCP Keep Alive packets also use empty data
// field, but they are not accessible in Golang directly. All this means that
// we should get information about the client's connection state from an outer
// world, so we receive an 'isBroken' flag. I hope that some day developers of
// Go language will make something working.
func (cp *PoolOfClients) TakeIdleClient(clientId string, isBroken bool) (err error) {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	cli, idExists := cp.usedClients[clientId]
	if !idExists {
		// Client is not being used !
		return ae.ErrorFromFormattedText(Errf_ClientIsNotBeingUsed, cli.GetId())
	}

	if isBroken {
		cp.brokenClients <- cli
	} else {
		cp.idleClients <- cli
	}
	delete(cp.usedClients, clientId)

	return nil
}
