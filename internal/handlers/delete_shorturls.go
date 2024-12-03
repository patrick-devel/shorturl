package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestDeleteShortURL []string

type serviceDeleter interface {
	DeleteShortURL(ctx context.Context, shortUrls []string) error
}

func DeleteShortUrls(service serviceDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqShortURLs RequestDeleteShortURL

		if err := c.BindJSON(&reqShortURLs); err != nil {
			c.JSON(http.StatusBadRequest, "")

			return
		}
		logrus.Infof("received: %v", reqShortURLs)

		go func() {
			err := service.DeleteShortURL(c.Copy(), reqShortURLs)
			if err != nil {
				logrus.WithError(err).Error("error deleting short urls")
			}
		}()

		c.JSON(http.StatusAccepted, http.NoBody)
	}
}
