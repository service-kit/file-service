package http

import (
	"github.com/gin-gonic/gin"
	"github.com/service-kit/file-service/config"
	"github.com/service-kit/file-service/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"sync"
)

type safeFilesystem struct {
	fs               http.FileSystem
	readDirBatchSize int
}

func (fs safeFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return StatFile{File: f, readDirBatchSize: fs.readDirBatchSize}, nil
}

type StatFile struct {
	http.File
	readDirBatchSize int
}

func (e StatFile) Stat() (os.FileInfo, error) {
	s, err := e.File.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
	LOOP:
		for {
			_, err := e.File.Readdir(e.readDirBatchSize)
			switch err {
			case io.EOF:
				break LOOP
			case nil:
				continue
			default:
				return nil, err
			}
		}
		return nil, os.ErrNotExist
	}
	return s, err
}

type httpManager struct {
	addr           string
	downAddr       string
	wg             *sync.WaitGroup
	engine         *gin.Engine
	authEnable     bool
	authPassEnable bool
}

var m *httpManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *httpManager {
	once.Do(func() {
		m = &httpManager{}
	})
	return m
}

func (self *httpManager) InitManager(wg *sync.WaitGroup) error {
	logger = log.GetInstance().GetLogger()
	self.wg = wg
	defer self.wg.Done()
	errCh := make(chan error, 1)
	var err error = nil
	self.engine = gin.Default()
	self.addr, err = config.GetInstance().GetConfig("HTTP_ADDR")
	if nil != err {
		return err
	}
	self.downAddr, err = config.GetInstance().GetConfig("FILE_DOWNLOAD_ADDR")
	if nil != err {
		return err
	}
	rootPath, err := config.GetInstance().GetConfig("FILE_ROOT_PATH")
	if nil != err {
		return err
	}
	self.engine.Static("/", "./public")
	self.engine.POST("/upload", handleFileUpload)
	go func() {
		errCh <- self.engine.Run(self.addr)
	}()
	fs := safeFilesystem{fs: http.Dir("./" + rootPath + "/"), readDirBatchSize: 2}
	fss := http.FileServer(fs)
	http.Handle("/", fss)
	go func() {
		errCh <- http.ListenAndServe(self.downAddr, nil)
	}()
	select {
	case err = <-errCh:
		logger.Error("http service shutdown", zap.Error(err))
		return err
	}
}
