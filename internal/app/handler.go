package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var urls = make(map[string]string)

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		shortener(w, r)
	case http.MethodGet:
		expander(w, r)
		t := ""
		_ = t
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func shortener(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Cannot read request body"))
		return
	}
	param := string(b)

	if _, err := url.Parse(param); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad URL given"))
		return
	}

	id := fmt.Sprintf("%x", len(urls))
	urls[id] = param

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("http://localhost:8080/"+ id))
}

func expander(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad ID given"))
		return
	}
	param := r.URL.Path[1:]

	target, ok := urls[param]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusTemporaryRedirect)
}