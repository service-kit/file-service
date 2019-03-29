package config

import (
	"errors"
	"github.com/service-kit/file-service/util"
	"github.com/Unknwon/goconfig"
	"io"
	"time"
)

type ServiceConfig struct {
	conf            *goconfig.ConfigFile
	confName        string
	isLoadSucc      bool
	fileLastModTime int64
	lastCheckTime   int64
}

func (self *ServiceConfig) Init(confFile string) error {
	self.confName = confFile
	var err error = nil
	self.conf, err = goconfig.LoadConfigFile(confFile)
	if nil != err {
		return err
	}
	self.isLoadSucc = true
	self.fileLastModTime, _ = util.GetFileModTime(self.confName)
	self.checkConfigFile()
	return err
}

func (self *ServiceConfig) InitFromData(data []byte) error {
	var err error = nil
	self.conf, err = goconfig.LoadFromData(data)
	if nil != err {
		return err
	}
	self.isLoadSucc = true
	return err
}

func (self *ServiceConfig) InitFromReader(reader io.Reader) error {
	var err error = nil
	self.conf, err = goconfig.LoadFromReader(reader)
	if nil != err {
		return err
	}
	self.isLoadSucc = true
	return err
}

func (self ServiceConfig) GetConfig(confName string) (string, error) {
	if !self.isLoadSucc {
		return "", errors.New("get config fail,is not load success!!!")
	}
	return self.conf.GetValue("", confName)
}

func (self *ServiceConfig) checkConfigFile() {
	go func() {
		for {
			fileModTime, err := util.GetFileModTime(self.confName)
			if nil != err {
				continue
			}
			if fileModTime != self.fileLastModTime {
				self.conf.Reload()
				self.fileLastModTime = fileModTime
				logger.Info(self.confName + " is change, conf reload")
			}
			time.Sleep(time.Second)
		}
	}()
}
