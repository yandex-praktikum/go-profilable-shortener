package main

import (
	"net/http"

	"github.com/bbrodriges/practicum-shortener/internal/app"
)

func main() {
	if err := run(); err != nil {
		panic("unexpected error: " + err.Error())
	}
}

func run() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(app.Router))
}
