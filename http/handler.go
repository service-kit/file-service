package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/service-kit/file-service/common"
	"github.com/service-kit/file-service/file"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func handleFileUpload(c *gin.Context) {
	ret := new(common.FSUploadResult)
	ret.Urls = make(map[string]string)
	ret.Code = 0
	statusCode := http.StatusOK
	defer func() {
		if err := recover(); err != nil {
			logger.Error("handle asr request error", zap.Reflect("error", err))
			ret.Code = -1
			ret.Msg = "service has some err"
			statusCode = http.StatusInternalServerError
		}
		logger.Info("handle upload request result", zap.String("req", c.Request.URL.RawQuery), zap.Any("ret", ret))
		c.JSON(statusCode, ret)
	}()

	defer func() {
		c.JSON(statusCode, ret)
	}()
	form, err := c.MultipartForm()
	if nil != err {
		ret.Code = -1
		ret.Msg = "not multipart"
		return
	}
	req := new(common.FSUploadParam)
	param := form.Value["param"]
	if nil == param {
		statusCode = http.StatusBadRequest
		ret.Code = -1
		ret.Msg = "not param"
		return
	}
	for i, p := range param {
		logger.Info("req param", zap.String(strconv.Itoa(i), p))
		json.Unmarshal([]byte(p), req)
	}
	now := time.Now().Unix()
	files := form.File["files"]
	for _, f := range files {
		fi := new(common.FileInfo)
		fi.Type = req.Type
		fi.Expire = req.Expire
		fi.Name = filepath.Base(f.Filename)
		fi.CreateTime = now
		in, _ := f.Open()
		fileAddr, err := file.GetInstance().SaveFile(fi, in)
		if nil != err {
			logger.Error("save file err", zap.Error(err))
		}
		ret.Urls[f.Filename] = fileAddr
	}
	return
}
