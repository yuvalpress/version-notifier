// version-notifier main
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"yuvalpress/version-notifier/internal/slack_notifier"

	jparser "github.com/Jeffail/gabs/v2"
	"github.com/Masterminds/semver/v3"
	xj "github.com/basgys/goxml2json"
	validate "golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
)

var (
	anchor Anchor

	// Reset color variables after call
	Reset = "\033[0m"

	// Green color for logs
	Green = "\033[32m"

	// Red color for logs
	Red = "\033[31m"
)

// Anchor holds the first initialized information for the service
type Anchor struct {
	repoList []Latest
}

// Latest holds all the needed information for a repo instance
type Latest struct {
	User   string
	Repo   string
	Latest string
	URL    string
}

// Conf struct holds all the repositories to configure
type Conf struct {
	Repos []map[string]string
}

// Init method for main Anchor object
func (a *Anchor) init() bool {
	confData, err := readConfigFile()
	if err != nil {
		log.Fatalf("Failed during initialization process with the following error: %v", err)
	}

	for _, info := range confData.Repos {
		for username, repoName := range info {
			data, err := download(username, repoName)
			if err != nil {
				log.Fatalf("Failed during initialization process with the following error: %v", err)
			}

			if len(data) == 0 {
				return false
			}

			a.repoList = append(a.repoList, Latest{
				User:   username,
				Repo:   repoName,
				Latest: getLatestTag(data[0]),
				URL:    data[0].Path("link.-href").String()})
		}
	}

	return true
}

// findRegexVersion returns the version inside a tag
func findRegexVersion(version string) string {
	r, _ := regexp.Compile("(\\d{1}|\\d{2}|\\d{3}|\\d{4})[.](\\d{1}|\\d{2}|\\d{3}|\\d{4})[.](\\d{1}|\\d{2}|\\d{3}|\\d{4})")
	match := r.FindString(version)

	return match
}

// stringInSlice returns true if string in list
func stringInSlice(level string, list []string) bool {
	for _, value := range list {
		if value == level {
			return true
		}
	}
	return false
}

// levelsToNotify returns a list of the levels to notify the user about
func levelsToNotify() []string {
	levels, exists := os.LookupEnv("NOTIFY")
	if !exists {
		return []string{"major", "minor", "patch"}
	}

	levels = strings.TrimSpace(strings.ToLower(levels))
	if levels == "all" {
		return []string{"major", "minor", "patch"}
	}

	return strings.Split(levels, ",")
}

// getURL build the github url with the needed user and repo
func getURL(username, repoName string) string {
	return "https://github.com/" + username + "/" + repoName + "/tags.atom"
}

// getLatestTag receives the latest ID (tag) available in the .atom file
func getLatestTag(data *jparser.Container) string {
	log.Println("Latest tag: v" + findRegexVersion(data.Path("id").String()))
	return "v" + findRegexVersion(data.Path("id").String())
}

// getUpdateLevel returns the update level: Major, Minor, Patch
// no need to validate this are semantic version formatted as this portion of the code is executed only after a test
func getUpdateLevel(old, new string) string {
	oldVer := semver.MustParse(old)
	newVer := semver.MustParse(new)

	if oldVer.Major() < newVer.Major() {
		return "major"
	} else if oldVer.Minor() < newVer.Minor() {
		return "minor"
	} else if oldVer.Patch() < newVer.Patch() {
		return "patch"
	}

	return "not any"
}

// doesNewTagExist receives two versions and validates if a newer version is available while validating both
// are in semantic version format
func doesNewTagExist(old, new string, repo string) (bool, string) {
	// validate version are indeed in semver format

	if validate.IsValid(old) && validate.IsValid(new) {
		oldVer := semver.MustParse(old)
		newVer := semver.MustParse(new)

		if oldVer.LessThan(newVer) {
			return true, newVer.String()
		}

		return false, ""
	}

	log.Printf(Red+"Something went wrong while trying to parse latest version from %v"+Reset, repo)
	return false, ""
}

// readConfigFile reads the repositories to scrape from the configmap attached to the pod as volume
func readConfigFile() (Conf, error) {
	var configData Conf
	conf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return Conf{}, err
	}

	err = yaml.Unmarshal(conf, &configData)

	if err != nil {
		return Conf{}, err
	}

	return configData, nil
}

// download is responsible to fetch the latest data from the relative url
func download(username, repoName string) ([]*jparser.Container, error) {
	url := getURL(username, repoName)
	log.Println("Fetching latest tags from:", url)

	// perform the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// convert XML to json
	json, err := xj.Convert(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	// load json to container object
	parseJSON, err := jparser.ParseJSON(json.Bytes())
	if err != nil {
		return nil, err
	}

	tagsData := parseJSON.Path("feed.entry").Children()

	if len(tagsData) == 0 {
		return nil, errors.New("request returned with 0 tags listed")
	}

	return tagsData, nil
}

// notify is responsible for notifying a selected Slack channel.
// in the future, more methods will be added
func notify(user, repo, url, oldVer, newVer string) {
	slack_notifier.Notify(user, repo, url, oldVer, newVer, getUpdateLevel(oldVer, newVer))
}

// main
func main() {
	// initialize application data until successful
	log.Println("Starting application...")

	log.Println("Initializing latest tags for configured repositories")
	for !anchor.init() {
		log.Printf("Failed to initialize application because of some bad requests...trying again.")
		time.Sleep(5 * time.Second)
		anchor.init()
	}

	levels := levelsToNotify()
	log.Printf("Notifications will be sent for: %s\n", levels)

	log.Println("Done!")
	log.Println("-----------------------------------------------------")

	// loop to infinity
	for true {
		time.Sleep(3 * time.Second)
		for index, repoData := range anchor.repoList {
			latest, _ := download(repoData.User, repoData.Repo)
			if latest != nil {
				result, newVer := doesNewTagExist(repoData.Latest, getLatestTag(latest[0]), repoData.User+"/"+repoData.Repo)

				if result {
					// update data
					anchor.repoList[index].Latest = "v" + newVer
					anchor.repoList[index].URL = latest[0].Path("link.-href").String()

					if stringInSlice(getUpdateLevel(repoData.Latest, newVer), levels) {
						log.Printf(Green+"New %v version found for package %v/%v: %v\n"+Reset,
							getUpdateLevel(repoData.Latest, newVer), repoData.User, repoData.Repo, newVer)

						// notify slack_notifier channel
						notify(repoData.User, repoData.Repo, anchor.repoList[index].URL, repoData.Latest, newVer)
					}
				} else {
					log.Printf("No new version found for package %v/%v", repoData.User, repoData.Repo)
				}
			}
		}
	}

}
