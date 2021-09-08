package app

import (
	"github.com/bbrodriges/practicum-shortener/internal/store"
)

type Instance struct {
	baseURL string

	store store.Store
}

func NewInstance(baseURL string, storage store.Store) *Instance {
	return &Instance{
		baseURL: baseURL,
		store:   storage,
	}
}
