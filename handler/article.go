package handler

import (
	"auditlog/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
}

func InitArticle() ArticleHandler {
	return ArticleHandler{}
}

func (h ArticleHandler) FetchArticles(c echo.Context) (err error) {
	datas := []models.Article{
		{
			ID:    "1",
			Title: "Hello World!",
			Body:  "No! Hi World!",
		},
	}

	return c.JSON(http.StatusOK, datas)
}
