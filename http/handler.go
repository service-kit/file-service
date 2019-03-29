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
)

func handleFileUpload(c *gin.Context) {
	ret := new(common.FSUploadResult)
	ret.Code = 0
	statusCode := http.StatusOK
	defer func() {
		c.JSON(statusCode,ret)
	}()
	form,err := c.MultipartForm()
	if nil != err {
		ret.Code = -1
		ret.Msg = "not multipart"
		return
	}
	rsp := make([]string,0)
	files := form.File["files"]
	for _, f := range files {
		fi := new(common.FileInfo)
		fi.Type = c.Param("type")
		fi.Expire,_ = strconv.ParseInt(c.Param("expire"),10,64)
		fi.Name = filepath.Base(f.Filename)
		in,_ := f.Open()
		fileAddr,err := file.GetInstance().SaveFile(fi,in)
		if nil != err {
			logger.Error("save file err",zap.Error(err))
		}
		rsp = append(rsp, fileAddr)
	}
	json.NewEncoder(c.Writer).Encode(rsp)
	return
}