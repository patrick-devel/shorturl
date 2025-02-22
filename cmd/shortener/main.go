package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/handlers"
	middlewares "github.com/patrick-devel/shorturl/internal/middlwares"
	"github.com/patrick-devel/shorturl/internal/models"
	"github.com/patrick-devel/shorturl/internal/service"
	"github.com/patrick-devel/shorturl/internal/storage"
)

type storager interface {
	ReadEvent(ctx context.Context, hash string) (string, error)
	WriteEvent(ctx context.Context, event models.Event) error
	WriteEvents(_ context.Context, events []models.Event) error
	ReadEventsByCreatorID(ctx context.Context, userID string) ([]models.Event, error)
	SetDeleteByShortURL(shorts []string) error
}

func makeMigrate(dsn string) {
	m, err := migrate.New("file://./migrations", dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := m.Up(); err != nil {
		logrus.Error(err)
	}
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

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		databaseDSN = parsedFlags.DatabaseDSN
	}

	fileStorage := os.Getenv("FILE_STORAGE_PATH")
	if fileStorage == "" {
		fileStorage = parsedFlags.FilePath
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

	var store storager
	var db *sql.DB

	switch {
	case cfg.DatabaseDSN != "":
		db, err = sql.Open("postgres", cfg.DatabaseDSN)
		if err != nil {
			logrus.Errorf("db not connected: %v", err)
		}

		if err = db.Ping(); err != nil {
			logrus.Errorf("db not connected: %v", err)
		}

		defer db.Close()
		makeMigrate(cfg.DatabaseDSN)

		store = storage.NewDBStorage(db, 5*time.Second)
	case fileStorage != "":
		store, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			logrus.Errorf("error fs %v", err)
		}
		defer cfg.RemoveTemp()
	default:
		cache := map[string]string{}
		store = storage.NewMemoryStorage(cache)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jwtSecret := os.Getenv("JWT_SIGNING_KEY")
	shortService := service.New(&cfg.BaseURL, store, ctx)
	authMidlwr := middlewares.AuthMiddleware(jwtSecret, logger)

	mux := gin.New()
	mux.Use(loggingMdlwr)
	mux.Use(middlewares.GzipMiddleware())
	mux.POST("/", authMidlwr, handlers.MakeShortLinkHandler(shortService))
	mux.GET(fmt.Sprintf("%s/:id", cfg.BaseURL.Path), handlers.RedirectShortLinkHandler(shortService))
	mux.POST("/api/shorten", authMidlwr, handlers.MakeShortURLJSONHandler(shortService))
	mux.POST("/api/shorten/batch", authMidlwr, handlers.MakeShortURLBulk(shortService))
	mux.GET("/api/user/urls", authMidlwr, handlers.GetURLsByCreatorID(shortService))
	mux.DELETE("/api/user/urls", authMidlwr, handlers.DeleteShortUrls(shortService))

	mux.GET("/ping", func(c *gin.Context) {
		if db != nil {
			if err := db.Ping(); err != nil {
				c.JSON(http.StatusInternalServerError, "")

				return
			}
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
