package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/bbrodriges/practicum-shortener/models"
)

func (i *Instance) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Cannot read request body"))
		return
	}

	shortURL, err := i.shorten(string(b))
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

	shortURL, err := i.shorten(req.URL)
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
	param := chi.URLParam(r, "id")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad ID given"))
		return
	}

	target, ok := i.urls[param]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (i *Instance) shorten(rawURL string) (shortURL string, err error) {
	if _, err := url.Parse(rawURL); err != nil {
		return "", errors.New("bad URL given")
	}

	id := fmt.Sprintf("%x", len(i.urls))
	i.urls[id] = rawURL
	return i.baseURL + id, nil
}