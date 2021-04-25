package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/core/log"
	"onesite/core/middleware"
	"onesite/core/model"
	"onesite/core/server/http/rest"
)

type FileResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Owner    string `json:"owner"`
	Exported bool   `json:"exported"`
}

func FileResponseFromUserModel(file *model.File) *FileResponse {
	return &FileResponse{
		Id:       file.ID,
		Name:     file.Name,
		Size:     file.Size,
		Owner:    file.Owner,
		Exported: file.Exported,
	}
}

func FileResponseFromUserModels(files []model.File) []*FileResponse {
	filesResponse := make([]*FileResponse, 0, len(files))
	for index := range files {
		filesResponse = append(filesResponse, FileResponseFromUserModel(&files[index]))
	}
	return filesResponse
}

type SetExportRequest struct {
	Exported bool `json:"exported" form:"exported"`
}

// ListFiles 分页查询文件列表
func (s *Service) ListFiles() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		count, files, err := s.Dao.ListFiles(
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
func (s *Service) DownloadFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		file, err := s.Dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		if file.Owner != user.Username {
			rest.PermissionDenied(c)
			return
		}

		response, err := s.Dao.Weed.DownloadFile(file.Fid, "")
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
func (s *Service) UploadFile() func(c *gin.Context) {
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

		fileInfo, err := s.Dao.Weed.UploadFile(src, filenamePart)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		fileIns, err := s.Dao.CreateFile(&model.File{
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
func (s *Service) DeleteFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		file, err := s.Dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		if file.Owner != user.Username {
			rest.PermissionDenied(c)
			return
		}

		tx := s.Dao.Orm.Db.Begin()
		err = s.Dao.DeleteFileWithDb(tx, request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		err = s.Dao.Weed.DeleteFile(file.Fid)
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
func (s *Service) SetExportFile() func(c *gin.Context) {
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
		file, err := s.Dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		if file.Owner != user.Username {
			rest.PermissionDenied(c)
			return
		}

		err = s.Dao.SetExportFile(request.PK, exportRequest.Exported)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		rest.NoContent(c)
	}
}

// ExportFile 外链访问文件
func (s *Service) ExportFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		file, err := s.Dao.QueryFile(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		if !file.Exported {
			rest.PermissionDenied(c)
			return
		}

		response, err := s.Dao.Weed.DownloadFile(file.Fid, "")
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
