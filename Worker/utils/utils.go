package utils

import "os"

func GetEnvValueWithDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
