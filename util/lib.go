package util

import (
	"github.com/subosito/gotenv"
	"os"
	"strings"
)

// Load .env variables
func init() {
	gotenv.Load()
}

// GetEnv takes in a key and returns value of the environment variable
func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); len(strings.TrimSpace(val)) > 0 {
		return val
	}
	return fallback
}
