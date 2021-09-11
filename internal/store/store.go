package store

import (
	"context"
	"io"
	"net/url"

	"github.com/gofrs/uuid"
)

type Store interface {
	io.Closer

	Save(ctx context.Context, url *url.URL) (id string, err error)
	Load(ctx context.Context, id string) (url *url.URL, err error)
}

type AuthStore interface {
	Store

	SaveUser(ctx context.Context, uid uuid.UUID, url *url.URL) (id string, err error)
	LoadUser(ctx context.Context, uid uuid.UUID, id string) (url *url.URL, err error)
	LoadUsers(ctx context.Context, uid uuid.UUID) (urls map[string]*url.URL, err error)
}
