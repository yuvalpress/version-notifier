package s3_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"sirrend/version-notifier/internal/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	client *s3.S3
}

// New creates a new instance of S3Client.
func New() (*S3Client, error) {
	profile, err := getAWSVariables()
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String(commons.REGION),
		},
	})
	if err != nil {
		log.Println("Error: Could not create a session, some configs are missing.")
		return nil, err
	}

	return &S3Client{client: s3.New(sess)}, nil
}

// getAWSVariables retrieves the AWS_PROFILE environment variable.
func getAWSVariables() (string, error) {
	profile, ok := os.LookupEnv("AWS_PROFILE")
	if !ok {
		return "", errors.New("AWS_PROFILE environment variable is not set")
	}
	return profile, nil
}

// Client returns the internal s3.S3 client.
func (s *S3Client) Client() *s3.S3 {
	return s.client
}

// Function to get an object from S3 bucket (stated in commons file)
// Input : fileName (This is the exact path inside the S3), outfile (The name and location of the file to save on the fileSystem)
func (c S3Client) getObject(fileName string, outputFile string) error {
	result, err := c.Client().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(commons.NOTIFIER_BUCKET),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer result.Body.Close()
	_, err = io.Copy(outFile, result.Body)
	if err != nil {
		return err
	}

	return nil
}

// Function to upload an object to S3 bucket (stated in commons file)
// Input : file (as interface), fileName (This is the exact path inside the S3)
func (c S3Client) UpdateObject(file interface{}, fileName string) error {
	jsonBytes, err := json.Marshal(file)
	if err != nil {
		log.Println("File failed to convert into a json")
		return err
	}

	fileToUpload := bytes.NewReader(jsonBytes)

	_, err = c.Client().PutObject(&s3.PutObjectInput{
		Bucket: aws.String(commons.NOTIFIER_BUCKET),
		Key:    aws.String(fileName),
		Body:   fileToUpload,
	})
	if err != nil {
		log.Println("Error: File failed to upload. Check if all values are configured properly.")
		return err
	}
	log.Println("File uploaded successfully")
	return nil
}

// Function to remove an object from an S3 bucket (stated in commons file)
// Input : fileName (This is the exact path inside the S3)
func (c S3Client) removeObject(fileName string) error {
	_, err := c.Client().DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(commons.NOTIFIER_BUCKET),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Println("Failed to delete file")
		return err
	}
	log.Println("File deleted successfully")
	return nil
}
