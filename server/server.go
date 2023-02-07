package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/vault-thirteen/Cache"
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/method"
	"github.com/vault-thirteen/SFRODB/common/protocol"
	"github.com/vault-thirteen/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/server/ff"
	"github.com/vault-thirteen/SFRODB/server/settings"
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

	mainListener net.Listener
	auxListener  net.Listener

	cacheT *cache.Cache[string, string]
	cacheB *cache.Cache[string, []byte]

	methodNameBuffers map[method.Method][]byte
	methodValues      map[string]method.Method

	filesT *ff.FilesFolder
	filesB *ff.FilesFolder

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

	srv.isRunning = new(atomic.Bool)
	srv.isRunning.Store(false)

	srv.cacheT = cache.NewCache[string, string](
		0,
		srv.settings.TextData.CacheVolumeMax,
		srv.settings.TextData.CachedItemTTL,
	)
	if err != nil {
		return nil, err
	}

	srv.cacheB = cache.NewCache[string, []byte](
		0,
		srv.settings.BinaryData.CacheVolumeMax,
		srv.settings.BinaryData.CachedItemTTL,
	)
	if err != nil {
		return nil, err
	}

	srv.filesT, err = ff.NewFilesFolder(srv.settings.TextData.Folder)
	if err != nil {
		return nil, err
	}

	srv.filesB, err = ff.NewFilesFolder(srv.settings.BinaryData.Folder)
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
	srv.mainListener, err = net.Listen(proto.LowLevelProtocol, srv.mainDsn)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	srv.auxListener, err = net.Listen(proto.LowLevelProtocol, srv.auxDsn)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
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

		conn, err := srv.mainListener.Accept()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
		} else {
			go srv.handleMainConnection(conn)
		}
	}

	log.Println("Main loop has stopped.")
}

func (srv *Server) runAuxLoop() {
	for {
		if !srv.isRunning.Load() {
			break
		}

		conn, err := srv.auxListener.Accept()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
		} else {
			go srv.handleAuxConnection(conn)
		}
	}

	log.Println("Auxiliary loop has stopped.")
}

// Stop stops the server.
func (srv *Server) Stop() (cerr *ce.CommonError) {
	var err error
	err = srv.mainListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = srv.auxListener.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	srv.isRunning.Store(false)
	// Main and Aux Loops will stop automatically.

	return nil
}

func (srv *Server) handleMainConnection(conn net.Conn) {
	c := connection.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
	)

	defer func() {
		derr := srv.finalize(c)
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *request.Request
	var err *ce.CommonError

	for {
		req, err = c.GetNextRequest()
		if err != nil {
			log.Println(err)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.ShowText,
			method.ShowBinary:
			err = srv.showRecord(c, req)

		case method.SearchTextRecord,
			method.SearchBinaryRecord:
			err = srv.searchRecord(c, req)

		case method.SearchTextFile,
			method.SearchBinaryFile:
			err = srv.searchFile(c, req)

		default:
			msg := fmt.Sprintf(ce.ErrUnsupportedMethodValue, req.Method)
			err = ce.NewClientError(msg, 0)
		}
		if err != nil {
			if err.IsServerError() {
				break
			} else {
				err = srv.clientError(c)
				if err != nil {
					break
				}
				continue
			}
		}
	}
}

func (srv *Server) handleAuxConnection(conn net.Conn) {
	c := connection.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
	)

	defer func() {
		derr := srv.finalize(c)
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *request.Request
	var err *ce.CommonError

	for {
		req, err = c.GetNextRequest()
		if err != nil {
			log.Println(err)
			break
		}

		if req.IsCloseConnection() {
			break
		}

		switch req.Method {
		case method.ForgetTextRecord,
			method.ForgetBinaryRecord:
			err = srv.forgetRecord(c, req)

		case method.ResetTextCache,
			method.ResetBinaryCache:
			err = srv.resetCache(c, req)

		default:
			msg := fmt.Sprintf(ce.ErrUnsupportedMethodValue, req.Method)
			err = ce.NewClientError(msg, 0)
		}
		if err != nil {
			if err.IsServerError() {
				break
			} else {
				err = srv.clientError(c)
				if err != nil {
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
func (srv *Server) finalize(c *connection.Connection) (cerr *ce.CommonError) {
	cerr = srv.closingConnection(c)
	if cerr != nil {
		return cerr
	}

	return c.Break()
}
