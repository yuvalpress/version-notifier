package commons

import "os"

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// exported program variables
var (
	// AWS
	NOTIFIER_BUCKET      = getEnvOrDefault("NOTIFIER_BUCKET", "")
	NOTIFIER_BUCKET_PATH = getEnvOrDefault("NOTIFIER_BUCKET", "version-notifier/")
	REGION               = getEnvOrDefault("REGION", "eu-west-1")
)
