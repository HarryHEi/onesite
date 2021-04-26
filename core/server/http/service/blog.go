package service

import (
	"github.com/gin-gonic/gin"

	"onesite/core/middleware"
	"onesite/core/model"
	"onesite/core/server/http/rest"
)

type CommitArticleRequest struct {
	Title    string `json:"title" form:"title" binding:"required,gte=3,lte=128"`
	Document string `json:"document" form:"document" binding:"required,gte=3"`
}

func (s *Service) CommitArticle() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request CommitArticleRequest
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user := middleware.ParseUser(c)
		err = s.Dao.CreateArticle(&model.Article{
			Author:   user.Username,
			Title:    request.Title,
			Document: request.Document,
			Comments: []model.Comment{},
		})
		rest.NoContent(c)
	}
}

func (s *Service) ListArticle() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		count, articles, err := s.Dao.QueryArticleView(request.Page, request.PageSize)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, rest.PaginationResponse{Count: count, Data: articles})
	}
}

type ArticlePictureRequest struct {
	PK string `uri:"pk"`
}

func (s *Service) ArticleDetail() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request ArticlePictureRequest
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		article, err := s.Dao.QueryArticleDetail(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, article)
		return
	}
}
