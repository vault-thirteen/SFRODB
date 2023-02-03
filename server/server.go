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
)

type Server struct {
	settings *settings.Settings
	dsn      string

	cacheT *cache.Cache[string, string]
	cacheB *cache.Cache[string, []byte]

	listener          net.Listener
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

	dsn := fmt.Sprintf("%s:%d", stn.ServerHost, stn.ServerPort)

	srv = &Server{
		settings: stn,
		dsn:      dsn,
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

func (srv *Server) GetDsn() (dsn string) {
	return srv.dsn
}

func (srv *Server) Start() (err error) {
	srv.listener, err = net.Listen(common.LowLevelProtocol, srv.dsn)
	if err != nil {
		return err
	}

	go srv.run()

	return nil
}

func (srv *Server) run() {
	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			log.Println(ErrConnectionAccepting, err.Error())
		}

		go srv.handleConnection(conn)
	}
}

func (srv *Server) Stop() (err error) {
	err = srv.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) handleConnection(conn net.Conn) {
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
