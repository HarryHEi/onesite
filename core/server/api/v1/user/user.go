package user

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/common/log"
	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/external"
	"onesite/core/middleware"
	"onesite/core/worker"
	"onesite/core/worker/tasks"
)

func Info() func(c *gin.Context) {
	return func(c *gin.Context) {
		user := middleware.ParseUser(c)
		rest.Success(c, InfoResponseFromUserModel(user))
	}
}

// UploadAvatar 用户上传头像
func UploadAvatar() func(c *gin.Context) {
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

		fileInfo, err := external.UploadFile(src, file.Filename)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)

		// 不必要的，仅用来测试的事件
		worker.ProduceTopic(tasks.DemoTopic, user.Username)

		if user.Avatar != "" {
			_ = external.DeleteFile(user.Avatar)
		}
		err = dao.UpdateUser(user.ID, map[string]interface{}{"avatar": fileInfo.Fid})
		if err != nil {
			_ = external.DeleteFile(fileInfo.Fid)
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}

// ExportAvatar 外链访问头像
func ExportAvatar() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user, err := dao.QueryUserById(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		if user.Avatar == "" {
			rest.NotFound(c)
			return
		}

		response, err := external.DownloadFile(user.Avatar, "height=200&width=200&mode=fill")
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
func UpdatePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request ChangePasswordRequest
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		user := middleware.ParseUser(c)
		err = dao.UpdateUser(
			user.ID,
			map[string]interface{}{
				"password": dao.GeneratePassword(request.Password),
			},
		)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
