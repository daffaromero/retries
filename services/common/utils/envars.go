package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	var failed []string

	envOS := getOSenv(key, &failed)
	envDotenv := getDotenv(key, &failed)

	if envOS != "" {
		return envOS
	} else if envDotenv != "" {
		return envDotenv
	}

	fmt.Printf("Failed to get vars from %s with key %s", failed, key)
	return ""
}

func getOSenv(key string, failed *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*failed = append(*failed, "os")
		return ""
	}

	return value
}

func getDotenv(key string, failed *[]string) string {
	err := godotenv.Load()
	if err != nil {
		*failed = append(*failed, ".env")
		return ""
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}

	return value
}
