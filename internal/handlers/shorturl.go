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
	MakeShortURL(ctx context.Context, originalURL, uid string) (string, error)
	GetOriginalURL(ctx context.Context, hash string) (string, error)
	MakeShortURLs(ctx context.Context, bulk models.ListRequestBulk) ([]models.Event, error)
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

		shortLink, err := service.MakeShortURL(c.Copy(), originalURL.String(), "")
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

		shortLink, err := service.MakeShortURL(c.Copy(), request.URL.String(), "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")

			return
		}

		c.JSON(http.StatusCreated, &models.Response{Result: shortLink})
	}
}

func MakeShortURLBulk(service shortService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.ListRequestBulk

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, "")

			return
		}

		if len(request) == 0 {
			c.JSON(http.StatusBadRequest, "")

			return
		}

		response := make(models.ListResponseBulk, 0, len(request))

		events, err := service.MakeShortURLs(c.Copy(), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "")

			return
		}

		for _, e := range events {
			response = append(response, models.ResponseBulk{
				ShortURL:      e.ShortURL,
				CorrelationID: e.UUID,
			})
		}

		c.JSON(http.StatusCreated, response)
	}
}
