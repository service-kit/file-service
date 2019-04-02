package file

import (
	"container/list"
	"errors"
	"github.com/service-kit/file-service/common"
	"github.com/service-kit/file-service/config"
	"github.com/service-kit/file-service/log"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
	"time"
)

type fileManager struct {
	localFileList  *list.List
	expireFileList *list.List
	mutex          sync.RWMutex
	basePath       string
	downDomain     string
}

var m *fileManager
var once sync.Once
var logger *zap.Logger

func GetInstance() *fileManager {
	once.Do(func() {
		m = &fileManager{}
	})
	return m
}

func (self *fileManager) InitManager() error {
	logger = log.GetInstance().GetLogger()
	self.localFileList = list.New()
	self.expireFileList = list.New()
	self.basePath, _ = config.GetInstance().GetConfig("FILE_ROOT_PATH")
	self.downDomain, _ = config.GetInstance().GetConfig("FILE_DOWNLOAD_DOMAIN")
	self.initFileDir()
	self.checkFile()
	return nil
}

func (self *fileManager) initFileDir() {
	err := os.MkdirAll(self.basePath+common.FP_AUDIO, 0777)
	if nil != err {
		logger.Error("mkdir err", zap.Error(err))
	}
	err = os.MkdirAll(self.basePath+common.FP_OTHER, 0777)
	if nil != err {
		logger.Error("mkdir err", zap.Error(err))
	}
}

func (self *fileManager) isFileExist(fn string) bool {
	if _, err := os.Stat(fn); nil != err {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (self *fileManager) SaveFile(fi *common.FileInfo, in io.Reader) (string, error) {
	if nil == fi || nil == in {
		return "", errors.New("file info is nil")
	}
	fileName := self.getFilePath(fi)
	if self.isFileExist(fileName) {
		return self.downDomain + fileName, nil
	}
	f, err := os.Create(self.basePath + fileName)
	if nil != err {
		return "", err
	}
	defer f.Close()
	io.Copy(f, in)
	self.mutex.Lock()
	self.localFileList.PushBack(fi)
	self.mutex.Unlock()
	return self.downDomain + fileName, nil
}

func (self *fileManager) checkFile() {
	delCh := make(chan *common.FileInfo, 256)
	delTicker := time.NewTicker(time.Second)
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now().Unix()
			self.mutex.Lock()
			for f := self.localFileList.Front(); nil != f; {
				fi := f.Value.(*common.FileInfo)
				if !fi.IsExpire(now) {
					break
				}
				if nil != f.Next() {
					f = f.Next()
					self.localFileList.Remove(f.Prev())
				} else {
					self.localFileList.Remove(f)
					f = f.Next()
				}
				delCh <- fi
			}
			self.mutex.Unlock()
		}
	}()
	go func() {
		for {
			select {
			case fi := <-delCh:
				self.expireFileList.PushBack(fi)
			case <-delTicker.C:
				for f := self.expireFileList.Front(); nil != f; {
					fi := f.Value.(*common.FileInfo)
					fileName := self.getFilePath(fi)
					if self.isFileExist(fileName) {
						os.Remove(fileName)
					}
					if nil != f.Next() {
						f = f.Next()
						self.expireFileList.Remove(f.Prev())
					} else {
						self.expireFileList.Remove(f)
						break
					}
				}
				delTicker = time.NewTicker(time.Second)
			}
		}
	}()
}

func (self *fileManager) getFilePath(fi *common.FileInfo) string {
	switch fi.Type {
	case common.FT_AUDIO_MP3, common.FT_AUDIO_WAV, common.FT_AUDIO_OPUS, common.FT_AUDIO_PCM:
		return common.FP_AUDIO + fi.Name
	default:
		return common.FP_OTHER + fi.Name
	}
}
