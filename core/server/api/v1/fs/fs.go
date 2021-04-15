package fs

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/common/config"
	"onesite/common/log"
	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/middleware"
)

// ListFiles 分页查询文件列表
func ListFiles() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		count, files, err := dao.ListFiles(
			[]string{"id", "name", "size", "owner"},
			request.Page,
			request.PageSize,
			"owner=?",
			user.Username,
		)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, rest.PaginationResponse{
			Count: count,
			Data:  FileResponseFromUserModels(files),
		})
	}
}

// DownloadFile 下载文件
func DownloadFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		file, err := dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		if file.Owner != user.Username {
			rest.PermissionDenied(c)
			return
		}

		// 从文件服务下载文件
		fileUrl := fmt.Sprintf(
			"%s://%s:%d/%s",
			config.CoreCfg.Weed.Protocol,
			config.CoreCfg.Weed.FsHost,
			config.CoreCfg.Weed.FsPort,
			file.Fid,
		)
		response, err := http.Get(fileUrl)
		if err != nil {
			rest.BadRequest(c, err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error("DownloadFile Body.Close()", zap.Error(err))
			}
		}(response.Body)
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", file.Name)
		extraHeaders := map[string]string{
			"Content-Disposition": contentDisposition,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, response.Body, extraHeaders)
	}
}

// UploadFile 上传文件
func UploadFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		file, _ := c.FormFile("file")

		src, err := file.Open()
		if err != nil {
			rest.BadRequest(c, err)
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				log.Error("UploadFile src.Close", zap.Error(err))
			}
		}(src)

		user := middleware.ParseUser(c)
		d, err := dao.UploadToFs(user.Username, src, file.Filename)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, d)
	}
}

// DeleteFile 删除文件
func DeleteFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		file, err := dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		if file.Owner != user.Username {
			rest.PermissionDenied(c)
			return
		}

		err = dao.DeleteFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		fileUrl := fmt.Sprintf(
			"%s://%s:%d/%s",
			config.CoreCfg.Weed.Protocol,
			config.CoreCfg.Weed.FsHost,
			config.CoreCfg.Weed.FsPort,
			file.Fid,
		)
		req, err := http.NewRequest(http.MethodDelete, fileUrl, nil)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
