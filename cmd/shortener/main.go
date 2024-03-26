package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/handlers"
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

	mux := gin.New()
	mux.POST("/", handlers.MakeShortLink(&cfg))
	mux.GET(fmt.Sprintf("%s/:id", cfg.BaseURL.Path), handlers.RedirectShortLink)
	mux.HandleMethodNotAllowed = true

	err = mux.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
