package server

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/vault-thirteen/Cache"
	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/server/ff"
	"github.com/vault-thirteen/SFRODB/server/settings"
)

const (
	ErrConnectionAccepting = "error accepting a connection: "
	MsgResettingCache      = "Resetting the Cache ..."
)

type Server struct {
	settings *settings.Settings

	mainDsn string
	auxDsn  string

	mainListener net.Listener
	auxListener  net.Listener

	cacheT *cache.Cache[string, string]
	cacheB *cache.Cache[string, []byte]

	methodNameBuffers map[common.Method][]byte
	methodValues      map[string]common.Method

	filesT *ff.FilesFolder
	filesB *ff.FilesFolder
}

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

	srv.methodNameBuffers, srv.methodValues = common.InitMethods()

	return srv, nil
}

func (srv *Server) GetMainDsn() (dsn string) {
	return srv.mainDsn
}

func (srv *Server) GetAuxDsn() (dsn string) {
	return srv.auxDsn
}

func (srv *Server) Start() (err error) {
	srv.mainListener, err = net.Listen(common.LowLevelProtocol, srv.mainDsn)
	if err != nil {
		return err
	}

	srv.auxListener, err = net.Listen(common.LowLevelProtocol, srv.auxDsn)
	if err != nil {
		return err
	}

	go srv.runMainLoop()
	go srv.runAuxLoop()

	return nil
}

func (srv *Server) runMainLoop() {
	for {
		conn, err := srv.mainListener.Accept()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
		}

		go srv.handleMainConnection(conn)
	}
}

func (srv *Server) runAuxLoop() {
	for {
		conn, err := srv.auxListener.Accept()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
		}

		go srv.handleAuxConnection(conn)
	}
}

func (srv *Server) Stop() (err error) {
	err = srv.mainListener.Close()
	if err != nil {
		return err
	}

	err = srv.auxListener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) handleMainConnection(conn net.Conn) {
	c, err := common.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
	)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		derr := c.Finalize()
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *common.Request

	for {
		req, err = c.GetNextRequest()
		if err != nil {
			log.Println(err)
			return
		}

		if req.IsCloseConnection() {
			return
		}

		switch req.Method {
		case common.MethodShowText:
			err = srv.showText(c, req)
		case common.MethodShowBinary:
			err = srv.showBinary(c, req)
		default:
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (srv *Server) handleAuxConnection(conn net.Conn) {
	c, err := common.NewConnection(
		conn,
		&srv.methodNameBuffers,
		&srv.methodValues,
		0,
	)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		derr := c.Finalize()
		if derr != nil {
			log.Println(derr)
		}
	}()

	var req *common.Request

	for {
		req, err = c.GetNextRequest()
		if err != nil {
			log.Println(err)
			return
		}

		if req.IsCloseConnection() {
			return
		}

		switch req.Method {
		case common.MethodForgetTextRecord:
			err = srv.forgetRecord(c, req)
		case common.MethodForgetBinaryRecord:
			err = srv.forgetRecord(c, req)
		case common.MethodResetTextCache:
			err = srv.resetCache(c, req)
		case common.MethodResetBinaryCache:
			err = srv.resetCache(c, req)
		default:
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (srv *Server) showText(c *common.Connection, r *common.Request) (err error) {
	var text string
	text, err = srv.getText(r.UID)
	if err != nil {
		return err
	}

	var rm *common.Response
	rm, err = common.NewResponse_ShowingText(text)
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) getText(uid string) (text string, err error) {
	// Check the UID.
	if !common.IsUidValid(uid) {
		return "", fmt.Errorf(common.ErrUid)
	}

	// Try to find the text in cache.
	text, err = srv.cacheT.GetRecord(uid)
	if err == nil {
		return text, nil
	}

	// Try the file storage.
	// Add an extension and convert path to the style of a current OS.
	relPath := filepath.Join(uid+srv.settings.TextData.FileExtension, "")
	var data []byte
	data, err = srv.filesT.GetFileContents(relPath)
	if err != nil {
		return "", err
	}
	text = string(data)

	// Save data in the cache.
	err = srv.cacheT.AddRecord(uid, text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (srv *Server) showBinary(c *common.Connection, r *common.Request) (err error) {
	var data []byte
	data, err = srv.getBinary(r.UID)
	if err != nil {
		return err
	}

	var rm *common.Response
	rm, err = common.NewResponse_ShowingBinary(data)
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) getBinary(uid string) (data []byte, err error) {
	// Check the UID.
	if !common.IsUidValid(uid) {
		return nil, fmt.Errorf(common.ErrUid)
	}

	// Try to find the data in cache.
	data, err = srv.cacheB.GetRecord(uid)
	if err == nil {
		return data, nil
	}

	// Try the file storage.
	// Add an extension and convert path to the style of a current OS.
	relPath := filepath.Join(uid+srv.settings.BinaryData.FileExtension, "")
	data, err = srv.filesB.GetFileContents(relPath)
	if err != nil {
		return nil, err
	}

	// Save data in the cache.
	err = srv.cacheB.AddRecord(uid, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (srv *Server) forgetRecord(c *common.Connection, r *common.Request) (err error) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return fmt.Errorf(common.ErrUid)
	}

	// Remove the record from the cache.
	switch r.Method {
	case common.MethodForgetTextRecord:
		err = srv.cacheT.RemoveRecord(r.UID)
		if err != nil {
			return err
		}

	case common.MethodForgetBinaryRecord:
		err = srv.cacheB.RemoveRecord(r.UID)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(common.ErrUnsupportedMethodValue, r.Method)
	}

	var rm *common.Response
	rm, err = common.NewResponse_OK()
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) resetCache(c *common.Connection, r *common.Request) (err error) {
	log.Println(MsgResettingCache)

	// Clear the cache.
	switch r.Method {
	case common.MethodResetTextCache:
		err = srv.cacheT.Clear()
		if err != nil {
			return err
		}

	case common.MethodResetBinaryCache:
		err = srv.cacheB.Clear()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf(common.ErrUnsupportedMethodValue, r.Method)
	}

	var rm *common.Response
	rm, err = common.NewResponse_OK()
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return nil
}
