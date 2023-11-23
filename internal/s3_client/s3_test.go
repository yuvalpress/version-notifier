package s3_client

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestS3FileUpload(t *testing.T) {
	fileName := "test_file.json"
	fileInterface := map[string]string{
		"name": "John Doe",
		"age":  "30",
	}
	jsonBytes, err := json.Marshal(fileInterface)
	if err != nil {
		log.Println("ERROR: File failed to convert into a json")
	}
	svc, err := New()
	if err != nil {
		t.Logf(err.Error())
	}

	err = svc.UpdateObject(jsonBytes, fileName)
	if err != nil {
		t.Errorf(err.Error())
		os.Exit(1)
	} else {
		svc.removeObject(fileName)
	}

	t.Logf("Test finished with success. Your function is working correctly :)")
}
