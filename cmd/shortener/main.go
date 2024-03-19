package main

import (
	"net/http"

	"github.com/patrick-devel/shorturl/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MakeShortLink)
	mux.HandleFunc("/{id}", handlers.RedirectShortLink)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
