package secrets_manager

import (
	"log"
	"os"
	"testing"
)

func TestSecretsManagerToVars(t *testing.T) {
	versionNotifierSecret, exists := os.LookupEnv("SECRET_NAME_TEST")
	if !exists {
		log.Println("INFO: Could not file SECRET_NAME_TEST as env var.")
		os.Exit(1)
	}

	err := ImportSecretsToEnv(versionNotifierSecret)
	if err != nil {
		t.Errorf(err.Error())
		os.Exit(1)
	}

	t.Logf("Test finished with success. Your function is working correctly :)")
}
