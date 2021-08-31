package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/bbrodriges/practicum-shortener/internal/app"
)

func newRouter(i *app.Instance) http.Handler {
	r := chi.NewRouter()

	r.Use(gzipMiddleware)
	r.Post("/", i.ShortenHandler)
	r.Post("/api/shorten", i.ShortenAPIHandler)
	r.Get("/{id}", i.ExpandHandler)

	return r
}

func gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
