package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	defaultAddr     = ":8080"
	defaultBaseURL  = url.URL{Scheme: "http", Host: "localhost:8080"}
	defaultLogLevel = "INFO"
)

type Config struct {
	Addr     string
	BaseURL  url.URL
	LogLevel string

	FileStoragePath string
	fileIsTemp      bool
}

func (c *Config) RemoveTemp() {
	if c.fileIsTemp {
		err := os.Remove(c.FileStoragePath)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{config: &Config{}}
}

func (cb *ConfigBuilder) WithAddress(addr string) *ConfigBuilder {
	if addr != ":0" {
		cb.config.Addr = addr
	}
	return cb
}

func (cb *ConfigBuilder) WithBaseURL(baseURL url.URL) *ConfigBuilder {
	if baseURL != (url.URL{}) {
		cb.config.BaseURL = baseURL
	}
	return cb
}

func (cb *ConfigBuilder) WithFileStoragePath(path string) *ConfigBuilder {
	if path != "" {
		cb.config.FileStoragePath = path
	}

	return cb
}

func (cb *ConfigBuilder) existOrCreateFile() error {
	_, err := os.Stat(cb.config.FileStoragePath)
	if errors.Is(err, os.ErrNotExist) {
		// если файла не сущ-ет, то создаем его
		_, err := os.Create(cb.config.FileStoragePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cb *ConfigBuilder) Build() (*Config, error) {
	if cb.config.Addr == "" || cb.config.Addr == ":0" {
		cb.config.Addr = defaultAddr
	}

	if cb.config.BaseURL == (url.URL{}) {
		cb.config.BaseURL = defaultBaseURL
	}

	if cb.config.FileStoragePath != "" {
		if err := cb.existOrCreateFile(); err != nil {
			return cb.config, fmt.Errorf("file path do not created: %w", err)
		}
		cb.config.fileIsTemp = false
	} else {
		file, err := os.CreateTemp("", "short-url-db.*.json")
		if err != nil {
			logrus.Fatal(err)
		}
		cb.config.FileStoragePath = file.Name()
		cb.config.fileIsTemp = true
	}

	cb.config.LogLevel = defaultLogLevel
	return cb.config, nil
}
