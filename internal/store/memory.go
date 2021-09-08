package store

import (
	"context"
	"fmt"
	"net/url"
)

type InMemory struct {
	store map[string]*url.URL
}

// NewInMemory create new InMemory instance
func NewInMemory() *InMemory {
	return &InMemory{
		store: make(map[string]*url.URL),
	}
}

func (m *InMemory) Save(_ context.Context, u *url.URL) (id string, err error) {
	id = fmt.Sprintf("%x", len(m.store))
	m.store[id] = u
	return id, nil
}

func (m *InMemory) Load(_ context.Context, id string) (u *url.URL, err error) {
	if u, ok := m.store[id]; ok {
		return u, nil
	}
	return nil, ErrNotFound
}

func (m *InMemory) Close() error {
	return nil
}
