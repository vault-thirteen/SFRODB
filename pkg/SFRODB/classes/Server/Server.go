package server

import (
	"net"
	"sync/atomic"

	"github.com/vault-thirteen/auxie/Cache/VL"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"

	client "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Connection"
	ff "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/FilesFolder"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Method"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Request"
	ss "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ServerSettings"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	Err_ConnectionAccepting = "error accepting a connection"
)

const (
	Msg_ResettingCache = "Resetting the Cache ..."
	Msg_MainLoopStop   = "Main loop has stopped."
	Msg_AuxLoopStop    = "Auxiliary loop has stopped."
)

// Server is server.
type Server struct {
	settings *ss.ServerSettings

	mainDsn string
	auxDsn  string

	mainListener     *net.TCPListener
	mainListenerAddr *net.TCPAddr

	auxListener     *net.TCPListener
	auxListenerAddr *net.TCPAddr

	cache *vl.Cache[string, []byte] // UID is string, Data is a byte array.
	files *ff.FilesFolder           // Data files.

	isRunning *atomic.Bool
}

func New(stn *ss.ServerSettings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings: stn,
		mainDsn:  common.FormatHostPort(stn.Hostname, stn.MainPort),
		auxDsn:   common.FormatHostPort(stn.Hostname, stn.AuxPort),
	}

	srv.mainListenerAddr, err = net.ResolveTCPAddr(protocol.LowLevelProtocol, srv.mainDsn)
	if err != nil {
		return nil, err
	}

	srv.auxListenerAddr, err = net.ResolveTCPAddr(protocol.LowLevelProtocol, srv.auxDsn)
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

	srv.files, err = ff.New(srv.settings.Data.Folder)
	if err != nil {
		return nil, err
	}

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
	srv.mainListener, err = net.ListenTCP(protocol.LowLevelProtocol, srv.mainListenerAddr)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, client.ClientId_None)
	}

	srv.auxListener, err = net.ListenTCP(protocol.LowLevelProtocol, srv.auxListenerAddr)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, client.ClientId_None)
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
			connErr := errors.ErrorWithPrefix(err, Err_ConnectionAccepting)
			srv.logError(connErr)
			continue
		}

		err = common.EnableTcpKeepAlives(conn, protocol.TcpKeepAliveIsEnabled, protocol.TcpKeepAlivePeriodSec)
		if err != nil {
			srv.logError(err)

			closeErr := conn.Close()
			if closeErr != nil {
				srv.logError(closeErr)
			}

			continue
		}

		go srv.handleMainConnection(conn)
	}

	srv.logMessage(Msg_MainLoopStop)
}

func (srv *Server) runAuxLoop() {
	for {
		if !srv.isRunning.Load() {
			break
		}

		conn, err := srv.auxListener.AcceptTCP()
		if err != nil {
			connErr := errors.ErrorWithPrefix(err, Err_ConnectionAccepting)
			srv.logError(connErr)
			continue
		}

		err = common.EnableTcpKeepAlives(conn, protocol.TcpKeepAliveIsEnabled, protocol.TcpKeepAlivePeriodSec)
		if err != nil {
			srv.logError(err)

			closeErr := conn.Close()
			if closeErr != nil {
				srv.logError(closeErr)
			}

			continue
		}

		go srv.handleAuxConnection(conn)
	}

	srv.logMessage(Msg_AuxLoopStop)
}

// Stop stops the server.
func (srv *Server) Stop() (cerr *ce.CommonError) {
	var err error
	err = srv.mainListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, client.ClientId_None)
	}

	err = srv.auxListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, client.ClientId_None)
	}

	srv.isRunning.Store(false)
	// Main and Aux Loops will stop automatically.

	return nil
}

func (srv *Server) handleMainConnection(conn *net.TCPConn) {
	con := connection.New(conn, 0, client.ClientId_Incoming)

	defer func() {
		derr := srv.finaliseConnection(con)
		if derr != nil {
			srv.logError(derr)
		}
	}()

	var req *request.Request
	var cerr *ce.CommonError

	for {
		req, cerr = con.GetNextRequest()
		if cerr != nil {
			srv.logError(cerr)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.Method_ShowData:
			cerr = srv.act_showData(con, req)
		case method.Method_SearchRecord:
			cerr = srv.act_searchRecord(con, req)
		case method.Method_SearchFile:
			cerr = srv.act_searchFile(con, req)
		default:
			cerr = ce.NewClientError(cio.FormatMessage(method.Errf_UnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
		}
		if cerr != nil {
			if cerr.IsServerError() {
				break
			} else {
				cerr = srv.respond_clientError(con)
				if cerr != nil {
					break
				}
				continue
			}
		}
	}
}

func (srv *Server) handleAuxConnection(conn *net.TCPConn) {
	con := connection.New(conn, 0, client.ClientId_Incoming)

	defer func() {
		derr := srv.finaliseConnection(con)
		if derr != nil {
			srv.logError(derr)
		}
	}()

	var req *request.Request
	var cerr *ce.CommonError

	for {
		req, cerr = con.GetNextRequest()
		if cerr != nil {
			srv.logError(cerr)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.Method_ForgetRecord:
			cerr = srv.act_forgetRecord(con, req)
		case method.Method_ResetCache:
			cerr = srv.act_resetCache(con, req)
		default:
			cerr = ce.NewClientError(cio.FormatMessage(method.Errf_UnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
		}
		if cerr != nil {
			if cerr.IsServerError() {
				break
			} else {
				cerr = srv.respond_clientError(con)
				if cerr != nil {
					break
				}
				continue
			}
		}
	}
}

// finaliseConnection is a method used by a Server to finalise the client's
// connection. This method is used either when the client requested to stop the
// communication or when an internal error happened on the server.
func (srv *Server) finaliseConnection(con *connection.Connection) (cerr *ce.CommonError) {
	cerr = srv.respond_closingConnection(con)
	if cerr != nil {
		return cerr
	}

	return con.Break()
}

func (srv *Server) logError(err error) {
	errors.LogError(err)
}

func (srv *Server) logMessage(msg string) {
	cio.LogMessage(msg)
}
