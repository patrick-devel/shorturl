package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	parsedFlags := ParseFlag()

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = parsedFlags.Addr.String()
	}

	baseURL, err := url.ParseRequestURI(os.Getenv("BASE_URL"))
	if baseURL != (&url.URL{}) || err != nil {
		baseURL = &parsedFlags.TemplateLink.url
	}

	fileStorage := os.Getenv("FILE_STORAGE_PATH")
	if fileStorage == "" {
		fileStorage = parsedFlags.FilePath
	}

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		databaseDSN = parsedFlags.DatabaseDSN
	}

	cfg, err := config.
		NewConfigBuilder().
		WithAddress(addr).
		WithBaseURL(*baseURL).
		WithFileStoragePath(fileStorage).
		WithDatabaseDSN(databaseDSN).
		Build()
	if err != nil {
		logrus.Fatal(fmt.Errorf("do not build config: %w", err))
	}

	defer cfg.RemoveTemp()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	loggingMdlwr := middlewares.LoggingMiddleware(logger)

	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		logrus.Fatalf("db not connected: %v", err)
	}

	defer db.Close()

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
	mux.GET("/ping", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, "")

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	mux.HandleMethodNotAllowed = true

	err = mux.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
