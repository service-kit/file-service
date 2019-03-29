package config

import (
	"github.com/service-kit/file-service/log"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
)

type configManager struct {
	conf ServiceConfig
}

var m *configManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *configManager {
	once.Do(func() {
		m = &configManager{}
	})
	return m
}

func (self *configManager) InitManager() error {
	logger = log.GetInstance().GetLogger()
	return self.conf.Init("./conf/file_service_conf.ini")
}

func (self configManager) GetConfig(confName string) (string, error) {
	return self.conf.GetConfig(confName)
}

func (self configManager) GetInt(confName string) (int, error) {
	value, err := self.conf.GetConfig(confName)
	if nil != err {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (self configManager) GetConfigArray(configName string) ([]string, error) {
	str, err := self.conf.GetConfig(configName)
	if nil != err {
		return nil, err
	}
	return strings.Split(str, ","), err
}
