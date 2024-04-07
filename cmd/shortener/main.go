package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/handlers"
	middlewares "github.com/patrick-devel/shorturl/internal/middlwares"
)

func main() {
	flagAddr, flagBaseURL := ParseFlag()

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = flagAddr.String()
	}

	baseURL, err := url.ParseRequestURI(os.Getenv("BASE_URL"))
	if baseURL != (&url.URL{}) || err != nil {
		baseURL = &flagBaseURL.url
	}

	cfg := config.NewConfigBuilder().WithAddress(addr).WithBaseURL(*baseURL).Build()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	loggingMdlwr := middlewares.LoggingMiddleware(logger)

	mux := gin.New()
	mux.Use(loggingMdlwr)
	mux.Use(middlewares.GzipMiddleware())
	mux.POST("/", handlers.MakeShortLinkHandler(&cfg))
	mux.GET(fmt.Sprintf("%s/:id", cfg.BaseURL.Path), handlers.RedirectShortLinkHandler)
	mux.POST("/api/shorten", handlers.MakeShortURLJSONHandler(&cfg))
	mux.HandleMethodNotAllowed = true

	err = mux.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
