package s3_client

import (
	"bytes"
	"io"
	"log"
	"os"

	"sirrend/internal/commons"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	client *s3.S3
}

// New creates a new instance of S3Client.
func New() (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(commons.REGION),
	})
	if err != nil {
		log.Println("Error: Could not create a session, some configs are missing.")
		return nil, err
	}

	return &S3Client{client: s3.New(sess)}, nil
}

// Client returns the internal s3.S3 client.
func (s *S3Client) Client() *s3.S3 {
	return s.client
}

// Function to get an object from S3 bucket (stated in commons file)
// Input : fileName (This is the exact path inside the S3), outfile (The name and location of the file to save on the fileSystem)
func (c S3Client) GetObject(fileName string) []byte {
	log.Println("INFO: the file name to fetch from bucket is: " + fileName)
	result, err := c.Client().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(commons.NOTIFIER_BUCKET),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("INFO: successfully fetched the file from S3.")
	defer result.Body.Close()

	// Read the YAML data
	yamlData, err := io.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return yamlData
}

// Function to upload an object to S3 bucket (stated in commons file)
// Input : file (as interface), fileName (This is the exact path inside the S3)
func (c S3Client) UpdateObject(data []byte, fileName string) error {
	fileToUpload := bytes.NewReader(data)

	_, err := c.Client().PutObject(&s3.PutObjectInput{
		Bucket: aws.String(commons.NOTIFIER_BUCKET),
		Key:    aws.String(fileName),
		Body:   fileToUpload,
	})
	if err != nil {
		log.Println("Error: File failed to upload. Check if all values are configured properly.")
		log.Println(err)
		return err
	}
	log.Println("INFO: File uploaded successfully!")
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
		log.Println("ERROR: Failed to delete file")
		return err
	}
	log.Println("INFO: File deleted successfully")
	return nil
}
