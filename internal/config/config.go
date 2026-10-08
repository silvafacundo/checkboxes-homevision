// Package config loads application configuration from the environment.
package config

import (
	"os"
	"strconv"
)

// Config holds runtime configuration for the application.
type Config struct {
	Port              string
	GinMode           string
	MaxUploadSize     int64
	AllowQueryOptions bool
}

// Load reads configuration from environment variables, applying defaults.
func Load() Config {
	return Config{
		Port:              getEnv("PORT", "8080"),
		GinMode:           getEnv("GIN_MODE", "debug"),
		MaxUploadSize:     getEnvInt64("MAX_UPLOAD_SIZE_MB", 20) << 20, // Value in MB, converted to bytes
		AllowQueryOptions: getEnvBool("ALLOW_QUERY_OPTIONS", false),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
	}
	return fallback
}
