package config

import (
	"flag"
	"os"
	"strings"
)

var (
	RunPort     = ":8080"
	BaseURL     = "http://localhost:8080/"
	PersistFile = ""
	AuthSecret  = []byte("ololo-trololo-shimba-boomba-look")
	DatabaseDSN = ""
)

func Parse() {
	flag.StringVar(&RunPort, "a", RunPort, "port to run server")
	flag.StringVar(&BaseURL, "b", BaseURL, "base URL for shorten URL response")
	flag.StringVar(&PersistFile, "f", PersistFile, "file to store shorten URLs")
	flag.StringVar(&DatabaseDSN, "d", DatabaseDSN, "connection string to database")

	flag.Parse()

	if val := os.Getenv("SERVER_ADDRESS"); val != "" {
		RunPort = val
	}
	if val := os.Getenv("BASE_URL"); val != "" {
		BaseURL = val
	}
	if val := os.Getenv("FILE_STORAGE_PATH"); val != "" {
		PersistFile = val
	}
	if val := os.Getenv("DATABASE_DSN"); val != "" {
		DatabaseDSN = val
	}

	BaseURL = strings.TrimRight(BaseURL, "/")
	_ = BaseURL
}
