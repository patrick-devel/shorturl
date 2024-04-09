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

var addr Addr
var template TemplateLink
var filePath string

func init() {
	flag.Var(&addr, "a", "Используйте адрес формата `host:port`")
	flag.Var(&template, "b", "Укажите адрес получения коротких ссылок. Пример: `http://localhost/path/to/short`")
	flag.StringVar(&filePath, "f", "", "Укажите путь до файла для сохранения данныз по запросам. Пример: `/path/to/dir`")
}

func ParseFlag() (Addr, TemplateLink, string) {
	flag.Parse()

	return addr, template, filePath
}
