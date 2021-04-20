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
	"onesite/core/external"
	"onesite/core/middleware"
	"onesite/core/model"
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
			[]string{"id", "name", "size", "owner", "exported"},
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

		var filenamePart string
		if len(file.Filename) > 32 {
			filenamePart = file.Filename[:32]
		} else {
			filenamePart = file.Filename
		}

		fileInfo, err := external.UploadFile(src, filenamePart)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		fileIns, err := dao.CreateFile(&model.File{
			Name:  filenamePart,
			Fid:   fileInfo.Fid,
			Size:  fileInfo.Size,
			Owner: user.Username,
		})
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, fileIns)
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

		daoIns, err := dao.GetDao()
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		tx := daoIns.Db.Begin()
		err = dao.DeleteFileWithDb(tx, request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		err = external.DeleteFile(file.Fid)
		if err != nil {
			tx.Rollback()
			rest.BadRequest(c, err)
			return
		}
		tx.Commit()
		rest.NoContent(c)
	}
}

// SetExportFile 设置文件为外链访问
func SetExportFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		var exportRequest SetExportRequest
		err = c.ShouldBind(&exportRequest)
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

		err = dao.SetExportFile(request.PK, exportRequest.Exported)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		rest.NoContent(c)
	}
}

// ExportFile 外链访问文件
func ExportFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		file, err := dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		if !file.Exported {
			rest.PermissionDenied(c)
			return
		}

		response, err := external.DownloadFile(file.Fid, "")
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error("DownloadFile Body.Close()", zap.Error(err))
			}
		}(response.Body)
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		contentDisposition := response.Header.Get("Content-Disposition")
		extraHeaders := map[string]string{
			"Content-Disposition": contentDisposition,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, response.Body, extraHeaders)
	}
}
