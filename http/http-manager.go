package http

import (
	"github.com/gin-gonic/gin"
	"github.com/service-kit/file-service/config"
	"github.com/service-kit/file-service/log"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

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
	self.engine.Static("/", "./public")
	self.engine.POST("/upload", handleFileUpload)
	go func() {
		errCh <- self.engine.Run(self.addr)
	}()
	http.Handle("/", http.FileServer(http.Dir("./file_root/")))
	go func() {
		errCh <- http.ListenAndServe(self.downAddr, nil)
	}()
	select {
	case err = <-errCh:
		logger.Error("http service shutdown", zap.Error(err))
		return err
	}
}
