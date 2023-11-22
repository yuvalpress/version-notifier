package s3_client

import (
	"os"
	"testing"
)

func TestS3FileUpload(t *testing.T) {
	fileName := "test_file.json"
	fileInterface := map[string]string{
		"name": "John Doe",
		"age":  "30",
	}
	svc, err := New()
	if err != nil {
		t.Logf(err.Error())
	}

	err = svc.UpdateObject(fileInterface, fileName)
	if err != nil {
		t.Errorf(err.Error())
		os.Exit(1)
	} else {
		svc.removeObject(fileName)
	}

	t.Logf("Test finished with success. Your function is working correctly :)")
}
