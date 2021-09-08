package main

import (
	"fmt"
	"net/http"

	"github.com/bbrodriges/practicum-shortener/internal/app"
	"github.com/bbrodriges/practicum-shortener/internal/config"
	"github.com/bbrodriges/practicum-shortener/internal/store"
)

func main() {
	if err := run(); err != nil {
		panic("unexpected error: " + err.Error())
	}
}

func run() error {
	config.Parse()

	var err error
	var storage store.Store = store.NewInMemory()
	if config.PersistFile != "" {
		storage, err = store.NewFileStore(config.PersistFile)
		if err != nil {
			return fmt.Errorf("cannot create persistent storage: %w", err)
		}
		defer storage.Close()
	}

	instance := app.NewInstance(config.BaseURL, storage)

	return http.ListenAndServe(config.RunPort, newRouter(instance))
}
