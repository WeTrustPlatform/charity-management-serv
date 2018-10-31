package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/subosito/gotenv"
)

// Load .env variables
func init() {
	if err := gotenv.Load(); err != nil {
		LogError(err)
	}
}

// GetEnv takes in a key and returns value of the environment variable
func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); len(strings.TrimSpace(val)) > 0 {
		return val
	}
	return fallback
}

// LogError handles errors gracefully
func LogError(err error) {
	fmt.Println(err)
}
