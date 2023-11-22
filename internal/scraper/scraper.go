package scraper

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	jparser "github.com/Jeffail/gabs/v2"
)

// GetURL build the needed GitHub urls with the received user and repo arguments
func GetURL(username, repoName string) (string, string) {
	return "https://api.github.com/repos/" + username + "/" + repoName + "/releases/latest",
		"https://api.github.com/repos/" + username + "/" + repoName + "/tags"
}

// getRequest returns a request with all the needed headers
func getRequest(url string) *http.Request {
	token, exist := os.LookupEnv("GITHUB_TOKEN")
	if !exist {
		log.Panicln("PANIC: The GITHUB_TOKEN environment variable must be set!")
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "version-notifier")

	return req
}

// request perform the actual request with the given url
func request(url, LogLevel string) (int, *http.Response) {
	// initialize request
	client := &http.Client{}
	req := getRequest(url)

	// perform the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil
	}

	return resp.StatusCode, resp
}

// bodyToJson receives a http.Response pointer object and returns parsed json as pointer
func bodyToJson(resp *http.Response) (*jparser.Container, error) {
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

// buildError build a readable error from the http.Response pointer object received
func buildError(resp *http.Response) error {
	r, _ := io.ReadAll(resp.Body)
	json, err := jparser.ParseJSON(r)
	if err != nil {
		return err
	}

	_ = resp.Body.Close()

	return errors.New(strconv.Itoa(resp.StatusCode) + ": request returned with the following message: " + json.Path("message").String())
}

func parseTagResponse(json *jparser.Container) *jparser.Container {
	return json.Children()[0]
}

// APIRequest uses a GitHub oauth token to retrieve needed data
func APIRequest(username, repoName, LogLevel string) (jsonObject *jparser.Container, requestType string, err error) {
	releaseURL, tagURL := GetURL(username, repoName)
	if LogLevel == "DEBUG" {
		log.Printf("Looking for tag and releases for %s/%s\n", username, repoName)
	}

	// initialize requests
	releaseRequestStatus, resp := request(releaseURL, LogLevel)

	// check if release is set for this repo
	if releaseRequestStatus == http.StatusOK {
		json, err := bodyToJson(resp)
		if err != nil {
			return nil, "", err
		}

		return json, "release", nil

	} else if releaseRequestStatus == http.StatusNotFound {
		tagRequestStatus, resp := request(tagURL, LogLevel)

		// check if tag is set for this repo
		if tagRequestStatus == http.StatusOK {
			json, err := bodyToJson(resp)
			if err != nil {
				return nil, "", err
			}

			return parseTagResponse(json), "tag", nil

		} else if tagRequestStatus == http.StatusNotFound {
			return nil, "", errors.New("No tag or release are set for " + username + "/" + repoName)
		}

	}

	return nil, "", buildError(resp)
}
