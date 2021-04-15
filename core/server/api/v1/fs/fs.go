package fs

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/common/log"
	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/middleware"
)

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
