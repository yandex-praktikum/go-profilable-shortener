package app

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadURLs(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test_urls.gob")
	defer os.Remove(path)

	fd, err := os.Create(path)
	require.NoError(t, err)

	defer fd.Close()

	content := map[string]string{"0x0": "https://praktikum.yandex.ru/"}

	enc := gob.NewEncoder(fd)
	err = enc.Encode(content)
	require.NoError(t, err)

	instance := &Instance{
		urls: make(map[string]string),
	}

	err = instance.LoadURLs(path)
	require.NoError(t, err)

	assert.Equal(t, content, instance.urls)
}

func TestStoreURLs(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test_urls.gob")
	defer os.Remove(path)

	fd, err := os.Create(path)
	require.NoError(t, err)

	defer fd.Close()

	instance := &Instance{
		urls: map[string]string{"0x0": "https://praktikum.yandex.ru/"},
	}

	err = instance.StoreURLs(path)
	require.NoError(t, err)

	_, err = fd.Seek(0, 0)
	require.NoError(t, err)

	var target map[string]string

	dec := gob.NewDecoder(fd)
	err = dec.Decode(&target)
	require.NoError(t, err)

	assert.Equal(t, instance.urls, target)
}
