package cp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/client"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/client/settings"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
)

const (
	ClientRestarterSuccessPauseSec = 5
	ClientRestarterFailurePauseSec = 15

	ErrDuplicateClientId       = "duplicate client id: %v"
	ErrNoIdleClientIsAvailable = "no idle client is available"
	ErrClientIsNotBeingUsed    = "client is not being used: %v"
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

func NewClientPool(size int, stn *settings.Settings) (pool *PoolOfClients, err error) {
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
		cli, err = client.NewClient(stn, strconv.Itoa(i))
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

	log.Printf("A pool of %d clients has been started.\r\n", cp.size)

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

	log.Println("Client restarter has stopped.")
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
	log.Println("Re-connecting the client ...")
	_ = cli.Stop()
	cerr := cli.Start()
	if cerr != nil {
		cp.brokenClients <- cli
		return false
	}

	log.Println("A broken client was successfully reconnected.")
	cp.idleClients <- cli
	return true
}

func (cp *PoolOfClients) Stop() {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	cp.mustStop.Store(true)

	// Drain all the channels except the 'stoppedClients' channel and stop all
	// the clients.
	log.Println("Stopping all the clients ...")
	var cli *client.Client
	for len(cp.idleClients) > 0 {
		cli = <-cp.idleClients
		_ = cli.Stop()
		log.Printf("Client [%s] was stopped.", cli.GetId())
		cp.stoppedClients <- cli
	}
	for i := range cp.usedClients {
		cli = cp.usedClients[i]
		_ = cli.Stop()
		log.Printf("Client [%s] was stopped.", cli.GetId())
		cp.stoppedClients <- cli
		delete(cp.usedClients, cli.GetId())
	}
	for len(cp.brokenClients) > 0 {
		cli = <-cp.brokenClients
		_ = cli.Stop()
		log.Printf("Client [%s] was stopped.", cli.GetId())
		cp.stoppedClients <- cli
	}

	log.Println("Waiting for subroutines to stop ...")
	cp.subRoutines.Wait()

	log.Println("Client pool has been stopped.")
}

// GiveIdleClient provides an idle client. If no clients are idle, an error is
// returned.
func (cp *PoolOfClients) GiveIdleClient() (cli *client.Client, err error) {
	cp.clientTransfers.Lock()
	defer cp.clientTransfers.Unlock()

	if len(cp.idleClients) == 0 {
		return nil, errors.New(ErrNoIdleClientIsAvailable)
	}

	cli = <-cp.idleClients
	_, idExists := cp.usedClients[cli.GetId()]
	if idExists {
		// Clients may not have same IDs !
		cp.idleClients <- cli
		return nil, fmt.Errorf(ErrDuplicateClientId, cli.GetId())
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
		return fmt.Errorf(ErrClientIsNotBeingUsed, cli.GetId())
	}

	if isBroken {
		cp.brokenClients <- cli
	} else {
		cp.idleClients <- cli
	}
	delete(cp.usedClients, clientId)

	return nil
}
