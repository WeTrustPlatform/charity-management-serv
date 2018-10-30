package main

import (
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); len(strings.TrimSpace(val)) > 0 {
		return val
	}
	return fallback
}
