package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	ss "github.com/vault-thirteen/SFRODB/pkg/SFHS/server/Settings"
	client "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	poc "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/PoolOfClients"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	ServerName = "SFHS"
)

const (
	HttpHdrCacheControlFormat = "max-age=%d, must-revalidate"
)

const (
	Msg_HttpErrorListenerStopped = "HTTP error listener has stopped."
	Msg_DbErrorListenerStopped   = "DB error listener has stopped."
)

const (
	Err_DbError = "DB error: "
)

type Server struct {
	settings *ss.Settings

	// HTTP(S) server.
	listenDsn  string
	httpServer *http.Server

	// DB client(s).
	dbDsnA        string
	dbDsnB        string
	poolOfClients *poc.PoolOfClients
	dbClientLock  *sync.Mutex

	// Channel for an external controller. When a message comes from this
	// channel, a controller must stop this server. The server does not stop
	// itself.
	mustBeStopped chan bool

	// Internal control structures.
	subRoutines *sync.WaitGroup
	mustStop    *atomic.Bool
	httpErrors  chan error
	dbErrors    chan *ce.CommonError

	// HTTP header values.
	httpHdrCacheControl string
}

func NewServer(stn *ss.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:      stn,
		listenDsn:     fmt.Sprintf(common.HostPortFormat, stn.ServerHost, stn.ServerPort),
		dbDsnA:        fmt.Sprintf(common.HostPortFormat, stn.DbHost, stn.DbPortA),
		dbDsnB:        fmt.Sprintf(common.HostPortFormat, stn.DbHost, stn.DbPortB),
		dbClientLock:  new(sync.Mutex),
		mustBeStopped: make(chan bool, 2),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, 8),
		dbErrors:      make(chan *ce.CommonError, 8),

		httpHdrCacheControl: fmt.Sprintf(HttpHdrCacheControlFormat, stn.HttpCacheControlMaxAge),
	}
	srv.mustStop.Store(false)

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	// DB Client.
	var dbClientSettings *cs.ClientSettings
	dbClientSettings, err = cs.New(
		srv.settings.DbHost,
		srv.settings.DbPortA,
		srv.settings.DbPortB,
		cs.ResponseMessageLengthLimitDefault,
	)
	if err != nil {
		return nil, err
	}

	srv.poolOfClients, err = poc.New(srv.settings.DbClientPoolSize, dbClientSettings)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
}

func (srv *Server) GetWorkMode() (modeId byte) {
	return srv.settings.ServerModeId
}

func (srv *Server) GetDbDsnA() (dsn string) {
	return srv.dbDsnA
}

func (srv *Server) GetDbDsnB() (dsn string) {
	return srv.dbDsnB
}

func (srv *Server) GetStopChannel() *chan bool {
	return &srv.mustBeStopped
}

func (srv *Server) Start() (cerr *ce.CommonError) {
	srv.startHttpServer()

	cerr = srv.poolOfClients.Start()
	if cerr != nil {
		return cerr
	}

	srv.subRoutines.Add(2)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()

	return nil
}

func (srv *Server) Stop() (cerr *ce.CommonError) {
	srv.mustStop.Store(true)

	var err error
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, client.ClientIdNone)
	}

	srv.poolOfClients.Stop()

	close(srv.httpErrors)
	close(srv.dbErrors)

	srv.subRoutines.Wait()

	return nil
}

func (srv *Server) startHttpServer() {
	go func() {
		var listenError error
		switch srv.settings.ServerModeId {
		case ss.ServerModeIdHttp:
			listenError = srv.httpServer.ListenAndServe()
		case ss.ServerModeIdHttps:
			listenError = srv.httpServer.ListenAndServeTLS(srv.settings.CertFile, srv.settings.KeyFile)
		}
		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) listenForHttpErrors() {
	defer srv.subRoutines.Done()

	for err := range srv.httpErrors {
		log.Println(common.Msg_ServerError + err.Error())
		srv.mustBeStopped <- true
	}

	log.Println(Msg_HttpErrorListenerStopped)
}

func (srv *Server) listenForDbErrors() {
	defer srv.subRoutines.Done()

	for cerr := range srv.dbErrors {
		log.Println(Err_DbError + cerr.Error())
	}

	log.Println(Msg_DbErrorListenerStopped)
}
