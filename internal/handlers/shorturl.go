package handlers

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/patrick-devel/shorturl/internal/models"
)

//go:generate mockgen -destination=./mocks/mock_service.go -package=mocks /Users/kspopova/GolandProjects/shorturl/internal/handlers shortService
type shortService interface {
	MakeShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, hash string) (string, error)
}

func MakeShortLinkHandler(service shortService) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())

			return
		}

		originalURL, err := url.ParseRequestURI(string(body))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())

			return
		}

		shortLink, err := service.MakeShortURL(c.Copy(), originalURL.String())
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.String(http.StatusCreated, shortLink)
	}
}

func RedirectShortLinkHandler(service shortService) gin.HandlerFunc {
	return func(c *gin.Context) {
		hashURL := c.Param("id")
		originalURL, err := service.GetOriginalURL(c.Copy(), hashURL)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())

			return
		}

		c.Redirect(http.StatusTemporaryRedirect, originalURL)
	}
}

func MakeShortURLJSONHandler(service shortService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.Request

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, "")

			return
		}

		shortLink, err := service.MakeShortURL(c.Copy(), request.URL.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")

			return
		}

		c.JSON(http.StatusCreated, &models.Response{Result: shortLink})
	}
}
