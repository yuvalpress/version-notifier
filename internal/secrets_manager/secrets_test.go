package secrets_manager

import (
	"os"
	"testing"
)

func TestSecretsManagerToVars(t *testing.T) {
	versionNotifierSecret, exists := os.LookupEnv("SECRET_NAME_TEST")
	if !exists {
		versionNotifierSecret = "SECRET"
	}

	err := importSecretsToEnv(versionNotifierSecret)
	if err != nil {
		t.Errorf(err.Error())
		os.Exit(1)
	}

	t.Logf("Test finished with success. Your function is working correctly :)")
}
