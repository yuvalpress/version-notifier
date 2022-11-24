package scraper

import (
	"errors"
	jparser "github.com/Jeffail/gabs/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// getRequest returns a request with all the needed headers
func getRequest(url string) *http.Request {
	token, exist := os.LookupEnv("GITHUB_TOKEN")
	if !exist {
		log.Panicln("The GITHUB_TOKEN environment variable must be set!")
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Barer "+token)
	req.Header.Set("User-Agent", "version-notifier")

	return req
}

// APIRequest uses a GitHub oauth token to retrieve needed data
func APIRequest(url, LogLevel string) (*jparser.Container, error) {
	if LogLevel == "DEBUG" {
		log.Println("Fetching latest release from:", url)
	}

	// initialize request
	client := &http.Client{}
	req := getRequest(url)

	// perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// convert to []byte
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		_ = resp.Body.Close()

		// load json to container object
		parseJSON, err := jparser.ParseJSON(bodyBytes)
		if err != nil {
			return nil, err
		}

		return parseJSON, nil
	}

	r, _ := io.ReadAll(resp.Body)
	json, err := jparser.ParseJSON(r)
	if err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return nil, errors.New(strconv.Itoa(resp.StatusCode) + ": request returned with the following message: " + json.Path("message").String())
}
