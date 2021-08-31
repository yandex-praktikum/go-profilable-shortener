package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bbrodriges/practicum-shortener/models"
)

func Test_run(t *testing.T) {
	go func() {
		err := run()
		require.NoError(t, err)
	}()

	time.Sleep(200 * time.Millisecond)

	targetURL := "https://praktikum.yandex.ru/"

	for i := 0; i < 50; i++ {
		expectedID := fmt.Sprintf("%x", i)

		t.Run("shorten", func(t *testing.T) {
			expectResponse := "http://localhost:8080/" + expectedID
			var actualResponse string

			{
				body := bytes.NewBufferString(targetURL)
				r := httptest.NewRequest("POST", "http://localhost:8080/", body)
				r.RequestURI = ""

				resp, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				require.Equal(t, http.StatusCreated, resp.StatusCode)

				defer resp.Body.Close()

				b, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				actualResponse = string(b)

				require.Equal(t, expectResponse, actualResponse)
			}

			{
				r := httptest.NewRequest("GET", actualResponse, nil)
				r.RequestURI = ""

				resp, err := http.DefaultTransport.RoundTrip(r)
				require.NoError(t, err)

				defer resp.Body.Close()

				assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
				assert.Equal(t, targetURL, resp.Header.Get("Location"))
			}
		})
	}

	for i := 50; i < 100; i++ {
		expectedID := fmt.Sprintf("%x", i)

		t.Run("shortenAPI", func(t *testing.T) {
			expectResponse := "{\"result\":\"http://localhost:8080/" + expectedID + "\"}\n"
			var actualResponse string

			{
				body := bytes.NewBufferString(`{"url":"` + targetURL + `"}`)
				r := httptest.NewRequest("POST", "http://localhost:8080/api/shorten", body)
				r.RequestURI = ""

				resp, err := http.DefaultClient.Do(r)
				require.NoError(t, err)
				require.Equal(t, http.StatusCreated, resp.StatusCode)

				defer resp.Body.Close()

				b, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				actualResponse = string(b)

				require.Equal(t, expectResponse, actualResponse)
			}

			{
				var target models.ShortenResponse
				err := json.Unmarshal([]byte(actualResponse), &target)
				require.NoError(t, err)

				r := httptest.NewRequest("GET", target.Result, nil)
				r.RequestURI = ""

				resp, err := http.DefaultTransport.RoundTrip(r)
				require.NoError(t, err)

				defer resp.Body.Close()

				assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
				assert.Equal(t, targetURL, resp.Header.Get("Location"))
			}
		})
	}

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(targetURL))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", "http://localhost:8080/", buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		expectResponse := "http://localhost:8080/64"
		actualResponse := string(b)

		require.Equal(t, expectResponse, actualResponse)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(targetURL)
		r := httptest.NewRequest("POST", "http://localhost:8080/", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		expectResponse := "http://localhost:8080/65"
		actualResponse := string(b)

		require.Equal(t, expectResponse, actualResponse)
	})
}
