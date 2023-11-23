package commons

import (
	"log"
	"os"
)

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Println("WARNING: Could not find the " + key + " environment variable. Using default value: " + defaultValue)
		return defaultValue
	}
	return value
}

// exported program variables
var (
	// AWS
	NOTIFIER_BUCKET      = getEnvOrDefault("NOTIFIER_BUCKET", "some_bucket")
	NOTIFIER_BUCKET_PATH = getEnvOrDefault("NOTIFIER_BUCKET_PATH", "version-notifier/")
	REGION               = getEnvOrDefault("REGION", "eu-west-1")
	CONFIG_FILE_NAME     = getEnvOrDefault("CONFIG_FILE_NAME", "config.yaml")
)
