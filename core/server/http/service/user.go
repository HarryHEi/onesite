package service

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/core/log"
	"onesite/core/middleware"
	"onesite/core/model"
	"onesite/core/server/http/rest"
	"onesite/core/worker/tasks"
)

type InfoResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `json:"avatar"`
}

func InfoResponseFromUserModel(user *model.User) *InfoResponse {
	return &InfoResponse{
		Id:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		IsAdmin:  user.IsAdmin,
		Avatar:   user.Avatar,
	}
}

type ChangePasswordRequest struct {
	Password string `json:"password" form:"password" binding:"required,gte=6,lte=32"`
}

func (s *Service) Info() func(c *gin.Context) {
	return func(c *gin.Context) {
		user := middleware.ParseUser(c)
		rest.Success(c, InfoResponseFromUserModel(user))
	}
}

// UploadAvatar 用户上传头像
func (s *Service) UploadAvatar() func(c *gin.Context) {
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

		fileInfo, err := s.Dao.Weed.UploadFile(src, file.Filename)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)

		// 不必要的，仅用来测试的事件
		s.Worker.ProduceTopic(tasks.DemoTopic, user.Username)

		if user.Avatar != "" {
			_ = s.Dao.Weed.DeleteFile(user.Avatar)
		}
		err = s.Dao.UpdateUser(user.ID, map[string]interface{}{"avatar": fileInfo.Fid})
		if err != nil {
			_ = s.Dao.Weed.DeleteFile(fileInfo.Fid)
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}

// ExportAvatar 外链访问头像
func (s *Service) ExportAvatar() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user, err := s.Dao.QueryUserById(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		if user.Avatar == "" {
			rest.NotFound(c)
			return
		}

		response, err := s.Dao.Weed.DownloadFile(user.Avatar, "height=200&width=200&mode=fill")
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

// UpdatePassword 更新密码
func (s *Service) UpdatePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request ChangePasswordRequest
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		user := middleware.ParseUser(c)
		err = s.Dao.UpdateUser(
			user.ID,
			map[string]interface{}{
				"password": s.Dao.GeneratePassword(request.Password),
			},
		)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
