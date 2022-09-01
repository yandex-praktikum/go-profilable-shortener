package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"

	"github.com/bbrodriges/practicum-shortener/internal/app"
	"github.com/bbrodriges/practicum-shortener/internal/auth"
)

func newRouter(i *app.Instance) http.Handler {
	r := chi.NewRouter()

	r.Use(authMiddleware)
	r.Post("/", i.ShortenHandler)
	r.Post("/api/shorten", i.ShortenAPIHandler)
	r.Post("/api/shorten/batch", i.BatchShortenAPIHandler)
	r.Delete("/api/user/urls", i.BatchRemoveAPIHandler)
	r.Get("/{id}", i.ExpandHandler)
	r.Get("/api/user/urls", i.UserURLsHandler)
	r.Get("/ping", i.PingHandler)

	return r
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var uid *uuid.UUID

		cookie, err := r.Cookie("auth")
		if cookie != nil {
			uid, err = auth.DecodeUIDFromHex(cookie.Value)
		}
		// generate new uid if failed to obtain existing
		if uid == nil {
			userID := ensureRandom()
			uid = &userID
		}

		// set new auth cookie in case of absence or decode error
		if err != nil {
			value, err := auth.EncodeUIDToHex(*uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("cannot encode auth cookie"))
				return
			}
			cookie = &http.Cookie{Name: "auth", Value: value}
			http.SetCookie(w, cookie)
		}

		// set uid to context
		ctx := auth.Context(r.Context(), *uid)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func ensureRandom() (res uuid.UUID) {
	for i := 0; i < 10; i++ {
		res = uuid.Must(uuid.NewV4())
	}
	return
}
