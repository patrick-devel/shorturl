package config

import "net/url"

var (
	defaultAddr    = ":8080"
	defaultBaseURL = url.URL{Scheme: "http", Host: "localhost:8080"}
)

type Config struct {
	Addr    string
	BaseURL url.URL
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

func (cb *ConfigBuilder) Build() Config {
	if cb.config.Addr == "" || cb.config.Addr == ":0" {
		cb.config.Addr = defaultAddr
	}

	if cb.config.BaseURL == (url.URL{}) {
		cb.config.BaseURL = defaultBaseURL
	}

	return *cb.config
}
