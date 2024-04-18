package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Addr struct {
	Host string
	Port int
}

func (a *Addr) Type() string {
	return a.String()
}

func (a *Addr) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *Addr) Set(value string) error {
	netAddr := strings.Split(value, ":")
	if len(netAddr) != 2 {
		return errors.New("need address in a form host:port")
	}

	port, err := strconv.Atoi(netAddr[1])
	if err != nil {
		return err
	}

	a.Host = netAddr[0]
	a.Port = port
	return nil
}

type TemplateLink struct {
	url url.URL
}

func (t *TemplateLink) Type() string {
	return t.String()
}

func (t *TemplateLink) String() string {
	return fmt.Sprint(*t)
}

func (t *TemplateLink) Set(value string) error {
	urlObj, err := url.ParseRequestURI(value)
	if err != nil {
		return err
	}
	t.url = *urlObj
	return nil
}

type ParsedFlags struct {
	Addr         Addr
	TemplateLink TemplateLink
	FilePath     string
	DatabaseDSN  string
}

var flags = &ParsedFlags{}

func init() {
	flag.Var(&flags.Addr, "a", "Используйте адрес формата `host:port`")
	flag.Var(&flags.TemplateLink, "b", "Адрес получения коротких ссылок. Пример: `http://localhost/path/to/short`")
	flag.StringVar(&flags.FilePath, "f", "/tmp/short-url-db.json", "Путь до файла для сохранения данных по запросам. Пример: `/path/to/dir`")
	flag.StringVar(&flags.DatabaseDSN, "d", "", "Адрес базы данных. Пример: postgresql://user:passwd@localhost:5432/dbname")
}

func ParseFlag() ParsedFlags {
	flag.Parse()

	return *flags
}
