package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/file_manager"
	"github.com/patrick-devel/shorturl/internal/handlers"
	middlewares "github.com/patrick-devel/shorturl/internal/middlwares"
	"github.com/patrick-devel/shorturl/internal/service"
)

func fileStorageSetup(logger *logrus.Logger, path string) *service.FileManager {
	consumer, err := filemanager.NewConsumer(path)
	if err != nil {
		logger.Fatal(err)
	}

	producer, err := filemanager.NewProducer(path)
	if err != nil {
		logger.Fatal(err)
	}

	return service.New(consumer, producer)
}

func main() {
	flagAddr, flagBaseURL, flagFilePath := ParseFlag()

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = flagAddr.String()
	}

	baseURL, err := url.ParseRequestURI(os.Getenv("BASE_URL"))
	if baseURL != (&url.URL{}) || err != nil {
		baseURL = &flagBaseURL.url
	}

	fileStorage := os.Getenv("FILE_STORAGE_PATH")
	if fileStorage == "" {
		fileStorage = flagFilePath
	}

	cfg, err := config.
		NewConfigBuilder().
		WithAddress(addr).
		WithBaseURL(*baseURL).
		WithFileStoragePath(fileStorage).
		Build()

	defer cfg.RemoveTemp()

	if err != nil {
		logrus.Fatal(fmt.Errorf("do not build config: %w", err))
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	loggingMdlwr := middlewares.LoggingMiddleware(logger)

	var fileManager *service.FileManager

	if fileStorage != "" {
		fileManager = fileStorageSetup(logger, fileStorage)
		defer fileManager.CloseFiles()
	}

	mux := gin.New()
	mux.Use(loggingMdlwr)
	mux.Use(middlewares.GzipMiddleware())
	mux.POST("/", handlers.MakeShortLinkHandler(cfg, fileManager))
	mux.GET(fmt.Sprintf("%s/:id", cfg.BaseURL.Path), handlers.RedirectShortLinkHandler(fileManager))
	mux.POST("/api/shorten", handlers.MakeShortURLJSONHandler(cfg, fileManager))
	mux.HandleMethodNotAllowed = true

	err = mux.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
