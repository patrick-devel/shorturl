package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/patrick-devel/shorturl/internal/handlers"
)

func main() {
	mux := gin.New()
	mux.POST("/", handlers.MakeShortLink)
	mux.GET("/:id", handlers.RedirectShortLink)
	mux.HandleMethodNotAllowed = true

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
