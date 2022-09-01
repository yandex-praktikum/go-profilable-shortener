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

	storage, err := newStore()
	if err != nil {
		return fmt.Errorf("cannot create storage: %w", err)
	}
	defer storage.Close()

	instance := app.NewInstance(config.BaseURL, storage)

	return http.ListenAndServe(config.RunPort, newRouter(instance))
}

func newStore() (storage store.AuthStore, err error) {
	return store.NewInMemory(), nil
}
