package config

import "net/url"

type Config struct {
	Addr         string
	TemplateLink url.URL
}

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{config: &Config{
		Addr:         ":8080",
		TemplateLink: url.URL{Scheme: "http", Host: "localhost:8080"}},
	}
}

func (cb *ConfigBuilder) WithAddress(addr string) *ConfigBuilder {
	if addr != ":0" {
		cb.config.Addr = addr
	}
	return cb
}

func (cb *ConfigBuilder) WithTemplateLink(templateLink url.URL) *ConfigBuilder {
	if templateLink != (url.URL{}) {
		cb.config.TemplateLink = templateLink
	}
	return cb
}

func (cb *ConfigBuilder) Build() Config {
	return *cb.config
}
