package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bbrodriges/practicum-shortener/internal/store"
	"github.com/bbrodriges/practicum-shortener/models"
)

func Test_ShortenAPIHandler(t *testing.T) {
	targetURL := "https://praktikum.yandex.ru/"

	instance := &Instance{
		baseURL: "http://localhost:8080",
		store:   store.NewInMemory(),
	}

	testCases := []struct {
		name             string
		url              string
		expectedStatus   int
		expectedResponse []byte
	}{
		{
			name:             "bad_request",
			url:              "htt_p://o.com",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: []byte("Cannot parse given string as URL"),
		},
		{
			name:             "success",
			url:              targetURL,
			expectedStatus:   http.StatusCreated,
			expectedResponse: []byte("{\"result\":\"http://localhost:8080/0\"}\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(models.ShortenRequest{URL: tc.url})
			require.NoError(t, err)
			body := bytes.NewBuffer(b)

			r := httptest.NewRequest("POST", "http://localhost:8080/api/shorten", body)
			w := httptest.NewRecorder()

			instance.ShortenAPIHandler(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedResponse, w.Body.Bytes())
		})
	}
}

func Test_expander(t *testing.T) {
	expectedURL := "https://praktikum.yandex.ru/"
	parsedURL, _ := url.Parse(expectedURL)

	storage := store.NewInMemory()
	id, _ := storage.Save(context.Background(), parsedURL)

	instance := &Instance{
		baseURL: "http://localhost:8080",
		store:   storage,
	}

	testCases := []struct {
		name             string
		id               string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "bad_request",
			id:               "",
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
		{
			name:             "not_found",
			id:               "-1",
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
		{
			name:             "success",
			id:               id,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: expectedURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://localhost:8080/"+tc.id, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			instance.ExpandHandler(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedLocation, w.Header().Get("Location"))
		})
	}
}
