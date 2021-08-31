package app

import (
	"encoding/gob"
	"os"
)

type Instance struct {
	baseURL string

	urls map[string]string
}

func NewInstance(baseURL string) *Instance {
	return &Instance{
		baseURL: baseURL,
		urls:    make(map[string]string),
	}
}

func (i *Instance) LoadURLs(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fd.Close()

	dec := gob.NewDecoder(fd)
	return dec.Decode(&i.urls)
}

func (i *Instance) StoreURLs(path string) error {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer fd.Close()

	enc := gob.NewEncoder(fd)
	return enc.Encode(i.urls)
}