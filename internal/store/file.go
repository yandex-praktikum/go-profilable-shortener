package store

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
)

type FileStore struct {
	hot     map[string]*url.URL
	enc     *gob.Encoder
	persist *os.File
}

// NewFileStore create new NewFileStore instance
func NewFileStore(filepath string) (*FileStore, error) {
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file at path %s: %w", filepath, err)
	}

	hot := make(map[string]*url.URL)

	dec := gob.NewDecoder(fd)
	if err := dec.Decode(&hot); err != nil {
		// truncate bad file for reuse
		if err := fd.Truncate(0); err != nil {
			return nil, fmt.Errorf("cannot truncate broken storage file: %w", err)
		}
	}

	return &FileStore{
		hot:     hot,
		enc:     gob.NewEncoder(fd),
		persist: fd,
	}, nil
}

func (f *FileStore) Save(_ context.Context, u *url.URL) (id string, err error) {
	id = fmt.Sprintf("%x", len(f.hot))
	f.hot[id] = u
	return id, f.flush()
}

func (f *FileStore) Load(_ context.Context, id string) (u *url.URL, err error) {
	if u, ok := f.hot[id]; ok {
		return u, nil
	}
	return nil, ErrNotFound
}

func (f *FileStore) Close() error {
	if err := f.flush(); err != nil {
		return fmt.Errorf("cannot flush data to file: %w", err)
	}
	return f.persist.Close()
}

func (f *FileStore) flush() error {
	return f.enc.Encode(f.hot)
}
