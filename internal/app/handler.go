package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/bbrodriges/practicum-shortener/internal/store"
	"github.com/bbrodriges/practicum-shortener/models"
)

func (i *Instance) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Cannot read request body"))
		return
	}

	u, err := url.Parse(string(b))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Cannot parse given string as URL"))
		return
	}

	shortURL, err := i.shorten(r.Context(), u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (i *Instance) ShortenAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Cannot parse given string as URL"))
		return
	}

	shortURL, err := i.shorten(r.Context(), u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(models.ShortenResponse{
		Result: shortURL,
	})

	if err != nil {
		fmt.Printf("cannot write response: %s", err)
	}
}

func (i *Instance) ExpandHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad ID given"))
		return
	}

	target, err := i.store.Load(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", target.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (i *Instance) shorten(ctx context.Context, rawURL *url.URL) (shortURL string, err error) {
	id, err := i.store.Save(ctx, rawURL)
	if err != nil {
		return "", fmt.Errorf("cannot save URL to storage: %w", err)
	}
	return i.baseURL + "/" + id, nil
}
