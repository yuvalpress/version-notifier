package secrets_manager

import (
	"encoding/json"
	"log"
	"os"

	"sirrend/internal/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func ImportSecretsToEnv(versionNotifierSecret string) error {
	// Create a new AWS session using the default credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(commons.REGION),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a Secrets Manager client
	secretsManagerClient := secretsmanager.New(sess)

	// Get the secret value from AWS Secrets Manager
	secretOutput, err := secretsManagerClient.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &versionNotifierSecret,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// Parse the secret value as a JSON or handle it based on your secret format
	secretData := *secretOutput.SecretString

	// Parse JSON string to map
	var secretMap map[string]string
	if err := json.Unmarshal([]byte(secretData), &secretMap); err != nil {
		log.Println(err)
		return err
	}

	// Set environment variables
	for key, value := range secretMap {
		os.Setenv(key, value)
	}

	log.Println("INFO: Secrets imported to environment variables successfully!")
	return nil
}
