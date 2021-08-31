package main

import (
	"net/http"

	"github.com/bbrodriges/practicum-shortener/internal/app"
	"github.com/bbrodriges/practicum-shortener/internal/config"
)

func main() {
	config.Parse()

	if err := run(); err != nil {
		panic("unexpected error: " + err.Error())
	}
}

func run() error {
	instance := app.NewInstance(config.BaseURL)

	if config.PersistFile != "" {
		_ = instance.LoadURLs(config.PersistFile)
		defer instance.StoreURLs(config.PersistFile)
	}

	return http.ListenAndServe(config.RunPort, newRouter(instance))
}
