package utils

import (
	"log"
	"os"
	"regexp"
	"sirrend/internal/scraper"
	"sirrend/internal/slack_notifier"
	"strconv"
	"strings"

	jparser "github.com/Jeffail/gabs/v2"
	"github.com/Masterminds/semver/v3"
	validate "golang.org/x/mod/semver"
)

var (
	// Reset color variables after call
	Reset = "\033[0m"

	// Red color for logs
	Red = "\033[31m"

	LogLevel = os.Getenv("LOG_LEVEL")
)

// GetUpdateLevel returns the update level: Major, Minor, Patch
// no need to validate this are semantic version formatted as this portion of the code is executed only after a test
func GetUpdateLevel(old, new string) string {
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

// FindRegexVersion returns the version inside a tag
func FindRegexVersion(version string) string {
	r, _ := regexp.Compile("(\\d{1}|\\d{2}|\\d{3}|\\d{4})[.](\\d{1}|\\d{2}|\\d{3}|\\d{4})[.](\\d{1}|\\d{2}|\\d{3}|\\d{4})")
	match := r.FindString(version)

	return match
}

// StringInSlice returns true if string in list
func StringInSlice(level string, list []string) bool {
	for _, value := range list {
		if strings.TrimSpace(value) == strings.TrimSpace(level) {
			return true
		}
	}
	return false
}

// DoesNewTagExist receives two versions and validates if a newer version is available while validating both
// are in semantic version format
func DoesNewTagExist(old, new string, repo string) (bool, string) {
	// validate version are indeed in semver format

	if validate.IsValid(old) && validate.IsValid(new) {
		oldVer := semver.MustParse(old)
		newVer := semver.MustParse(new)

		if oldVer.LessThan(newVer) {
			log.Println("INFO: Found a new version: " + newVer.String())
			return true, newVer.String()
		}

		return false, ""
	}

	log.Printf(Red+"INFO: Something went wrong while trying to parse latest version from %v"+Reset, repo)
	return false, ""
}

// GetLatestTag receives the latest ID (tag) available in the .atom file
func GetLatestTag(data, LogLevel string) string {
	if LogLevel == "DEBUG" {
		log.Println("DEBUG: Latest tag: v" + FindRegexVersion(data))
	}
	return "v" + FindRegexVersion(data)
}

// LevelsToNotify returns a list of the levels to notify the user about
func LevelsToNotify() []string {
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

// GetVersion is responsible to fetch the latest data from the relative url
func GetVersion(username, repoName string) (*jparser.Container, string, error) {
	json, requestType, err := scraper.APIRequest(username, repoName, LogLevel)
	if err != nil {
		return nil, "", err
	}

	return json, requestType, nil
}

// Notify is responsible for notifying a selected Slack channel.
// in the future, more methods will be added
func Notify(user, repo, url, oldVer, newVer, versionType string) {
	method, found := os.LookupEnv("NOTIFICATION_METHOD")
	if !found {
		log.Panicln("The NOTIFICATION_METHOD environment variable must be set!")
	}

	sendFullChangelog, found := os.LookupEnv("SEND_FULL_CHANGELOG")
	if !found {
		log.Println("INFO: The SEND_FULL_CHANGELOG environment variable is not set! Defaulting to `false`.")
		sendFullChangelog = "false"
	}

	// convert to bool
	sendBool, err := strconv.ParseBool(sendFullChangelog)
	if err != nil {
		log.Panicf("The SEND_FULL_CHANGELOG environment variable must be set to true or false only!")
	}

	if method == "none" {
		log.Panicln("The NOTIFICATION_METHOD environment variable must be set!")

	} else if method == "slack" {
		slack_notifier.Notify(user, repo, url, oldVer, newVer, GetUpdateLevel(oldVer, newVer), versionType, sendBool)
	}
}
