package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/vault-thirteen/Cache/VL"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/client"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method"
	proto "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/protocol"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/server/ff"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/server/settings"
)

const (
	ErrConnectionAccepting = "error accepting a connection: "
	MsgResettingCache      = "Resetting the Cache ..."
)

// Server is server.
type Server struct {
	settings *settings.Settings

	mainDsn string
	auxDsn  string

	mainListener     *net.TCPListener
	mainListenerAddr *net.TCPAddr

	auxListener     *net.TCPListener
	auxListenerAddr *net.TCPAddr

	methodNameBuffers map[method.Method][]byte
	methodValues      map[string]method.Method

	cache *vl.Cache[string, []byte] // UID is string, Data is a byte array.
	files *ff.FilesFolder           // Data files.

	isRunning *atomic.Bool
}

// NewServer creates a server.
func NewServer(stn *settings.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings: stn,
		mainDsn:  fmt.Sprintf("%s:%d", stn.ServerHost, stn.MainPort),
		auxDsn:   fmt.Sprintf("%s:%d", stn.ServerHost, stn.AuxPort),
	}

	srv.mainListenerAddr, err = net.ResolveTCPAddr(proto.LowLevelProtocol, srv.mainDsn)
	if err != nil {
		return nil, err
	}

	srv.auxListenerAddr, err = net.ResolveTCPAddr(proto.LowLevelProtocol, srv.auxDsn)
	if err != nil {
		return nil, err
	}

	srv.isRunning = new(atomic.Bool)
	srv.isRunning.Store(false)

	srv.cache = vl.NewCache[string, []byte](
		0,
		srv.settings.Data.CacheVolumeMax,
		srv.settings.Data.CachedItemTTL,
	)

	srv.files, err = ff.NewFilesFolder(srv.settings.Data.Folder)
	if err != nil {
		return nil, err
	}

	srv.methodNameBuffers, srv.methodValues = method.InitMethods()

	return srv, nil
}

// GetMainDsn returns the DSN of the main connection.
func (srv *Server) GetMainDsn() (dsn string) {
	return srv.mainDsn
}

// GetAuxDsn returns the DSN of the auxiliary connection.
func (srv *Server) GetAuxDsn() (dsn string) {
	return srv.auxDsn
}

// Start starts the server.
func (srv *Server) Start() (cerr *ce.CommonError) {
	var err error
	srv.mainListener, err = net.ListenTCP(proto.LowLevelProtocol, srv.mainListenerAddr)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, client.ClientIdNone)
	}

	srv.auxListener, err = net.ListenTCP(proto.LowLevelProtocol, srv.auxListenerAddr)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, client.ClientIdNone)
	}

	srv.isRunning.Store(true)
	go srv.runMainLoop()
	go srv.runAuxLoop()

	return nil
}

func (srv *Server) runMainLoop() {
	for {
		if !srv.isRunning.Load() {
			break
		}

		conn, err := srv.mainListener.AcceptTCP()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
			continue
		}

		err = connection.EnableKeepAlives(conn)
		if err != nil {
			log.Println(err.Error())
			closeErr := conn.Close()
			if closeErr != nil {
				log.Println(err.Error())
			}
			continue
		}

		go srv.handleMainConnection(conn)
	}

	log.Println("Main loop has stopped.")
}

func (srv *Server) runAuxLoop() {
	for {
		if !srv.isRunning.Load() {
			break
		}

		conn, err := srv.auxListener.AcceptTCP()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
			continue
		}

		err = connection.EnableKeepAlives(conn)
		if err != nil {
			log.Println(err.Error())
			closeErr := conn.Close()
			if closeErr != nil {
				log.Println(err.Error())
			}
			continue
		}

		go srv.handleAuxConnection(conn)
	}

	log.Println("Auxiliary loop has stopped.")
}

// Stop stops the server.
func (srv *Server) Stop() (cerr *ce.CommonError) {
	var err error
	err = srv.mainListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, client.ClientIdNone)
	}

	err = srv.auxListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, client.ClientIdNone)
	}

	srv.isRunning.Store(false)
	// Main and Aux Loops will stop automatically.

	return nil
}

func (srv *Server) handleMainConnection(conn *net.TCPConn) {
	con := connection.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
		client.ClientIdIncoming,
	)

	defer func() {
		derr := srv.finalize(con)
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *request.Request
	var cerr *ce.CommonError

	for {
		req, cerr = con.GetNextRequest()
		if cerr != nil {
			log.Println(cerr)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.ShowData:
			cerr = srv.showData(con, req)
		case method.SearchRecord:
			cerr = srv.searchRecord(con, req)
		case method.SearchFile:
			cerr = srv.searchFile(con, req)
		default:
			msg := fmt.Sprintf(ce.ErrUnsupportedMethodValue, req.Method)
			cerr = ce.NewClientError(msg, 0, con.ClientId())
		}
		if cerr != nil {
			if cerr.IsServerError() {
				break
			} else {
				cerr = srv.clientError(con)
				if cerr != nil {
					break
				}
				continue
			}
		}
	}
}

func (srv *Server) handleAuxConnection(conn *net.TCPConn) {
	con := connection.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
		client.ClientIdIncoming,
	)

	defer func() {
		derr := srv.finalize(con)
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *request.Request
	var cerr *ce.CommonError

	for {
		req, cerr = con.GetNextRequest()
		if cerr != nil {
			log.Println(cerr)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.ForgetRecord:
			cerr = srv.forgetRecord(con, req)
		case method.ResetCache:
			cerr = srv.resetCache(con, req)
		default:
			msg := fmt.Sprintf(ce.ErrUnsupportedMethodValue, req.Method)
			cerr = ce.NewClientError(msg, 0, con.ClientId())
		}
		if cerr != nil {
			if cerr.IsServerError() {
				break
			} else {
				cerr = srv.clientError(con)
				if cerr != nil {
					break
				}
				continue
			}
		}
	}
}

// finalize is a method used by a Server to finalize the client's connection.
// This method is used either when the client requested to stop the
// communication or when an internal error happened on the server.
func (srv *Server) finalize(con *connection.Connection) (cerr *ce.CommonError) {
	cerr = srv.closingConnection(con)
	if cerr != nil {
		return cerr
	}

	return con.Break()
}
