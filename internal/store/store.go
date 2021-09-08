package store

import (
	"context"
	"io"
	"net/url"
)

type Store interface {
	io.Closer

	Save(ctx context.Context, url *url.URL) (id string, err error)
	Load(ctx context.Context, id string) (url *url.URL, err error)
}
