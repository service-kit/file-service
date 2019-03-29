package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type logManager struct {
	logger *zap.Logger
}

var m *logManager
var once sync.Once

func GetInstance() *logManager {
	once.Do(func() {
		m = &logManager{}
	})
	return m
}

func (self *logManager) InitManager() (err error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	self.logger, err = cfg.Build()
	if nil != err {
		return
	}
	self.logger = self.logger.Named("file_service")
	return
}

func (self *logManager) FinishProcess() {
	if nil != self.logger {
		self.logger.Sync()
	}
}

func (self *logManager) GetLogger() *zap.Logger {
	return self.logger
}
