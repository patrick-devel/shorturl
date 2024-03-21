package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/patrick-devel/shorturl/config"
	"github.com/patrick-devel/shorturl/internal/handlers"
)

func main() {
	addr, template := ParseFlag()
	cfg := config.
		NewConfigBuilder().
		WithAddress(addr.String()).
		WithTemplateLink(template.url).
		Build()

	mux := gin.New()
	mux.POST("/", handlers.MakeShortLink(&cfg))
	mux.GET(fmt.Sprintf("%s/:id", cfg.TemplateLink.Path), handlers.RedirectShortLink)
	mux.HandleMethodNotAllowed = true

	err := mux.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
