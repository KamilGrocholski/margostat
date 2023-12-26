package config

import (
	"log"
	"os"
	"strconv"
)

var (
	DATABASE_URL                       string
	STATIC_DIR_PATH                    string
	SCRAPER_INTERVAL_IN_MINUTES        int64
	SCRAPER_DB_CONNECTION_MAX_ATTEMPTS int64
	SCRAPER_MAX_ATTEMPTS               int64
)

func Init() {
	DATABASE_URL = getStringEnv("DATABASE_URL")
	STATIC_DIR_PATH = getStringEnv("STATIC_DIR_PATH")
	SCRAPER_INTERVAL_IN_MINUTES = getInt64Env("SCRAPER_INTERVAL_IN_MINUTES")
	SCRAPER_DB_CONNECTION_MAX_ATTEMPTS = getInt64Env("SCRAPER_DB_CONNECTION_MAX_ATTEMPTS")
	SCRAPER_MAX_ATTEMPTS = getInt64Env("SCRAPER_MAX_ATTEMPTS")
}

func getStringEnv(name string) string {
	raw := os.Getenv(name)
	if raw == "" {
		log.Fatalf("Empty env %s", name)
	}

	return raw
}

func getInt64Env(name string) int64 {
	raw := getStringEnv(name)

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		log.Fatalf("Env %s is not int64", name)
	}

	return parsed
}
