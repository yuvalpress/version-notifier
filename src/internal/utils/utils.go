package utils

import (
	jparser "github.com/Jeffail/gabs/v2"
	"github.com/Masterminds/semver/v3"
	validate "golang.org/x/mod/semver"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yuvalpress/version-notifier/internal/scraper"
	"yuvalpress/version-notifier/internal/slack_notifier"
	"yuvalpress/version-notifier/internal/telegram_notifier"
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
			return true, newVer.String()
		}

		return false, ""
	}

	log.Printf(Red+"Something went wrong while trying to parse latest version from %v"+Reset, repo)
	return false, ""
}

// GetLatestTag receives the latest ID (tag) available in the .atom file
func GetLatestTag(data, LogLevel string) string {
	if LogLevel == "DEBUG" {
		log.Println("Latest tag: v" + FindRegexVersion(data))
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

// GetInterval returns the interval to use on run
func GetInterval() (string, int) {
	interval := os.Getenv("INTERVAL")
	intInterval, err := strconv.Atoi(interval)
	if err != nil {
		log.Println(Red + "Wrong INTERVAL environment variable inserted, defaulting to 30 min" + Reset)
		return "30", 30
	}

	return interval, intInterval
}

// WaitForInterval is responsible for waiting the interval requested by the user
func WaitForInterval() {
	strInterval, intInterval := GetInterval()
	log.Println("Performing next request in: " + strInterval + " minutes.")

	//wait <interval> minutes for next run
	time.Sleep(time.Duration(intInterval) * time.Minute)
	log.Println("Starting new run...")
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
func Notify(user, repo, url, oldVer, newVer string) {
	method, found := os.LookupEnv("NOTIFICATION_METHOD")
	if !found {
		log.Panicln("The NOTIFICATION_METHOD environment variable must be set!")
	}

	sendFullChangelog, found := os.LookupEnv("SEND_FULL_CHANGELOG")
	if !found {
		log.Println("The SEND_FULL_CHANGELOG environment variable is not set! Defaulting to `false`")
	}

	// convert to bool
	sendBool, err := strconv.ParseBool(sendFullChangelog)
	if err != nil {
		log.Panicf("The SEND_FULL_CHANGELOG environment variable must be set to true or false only!")
	}

	if method == "none" {
		log.Panicln("The NOTIFICATION_METHOD environment variable must be set!")

	} else if method == "telegram" {
		telegram_notifier.Notify(user, repo, url, oldVer, newVer, GetUpdateLevel(oldVer, newVer), sendBool)

	} else if method == "slack" {
		slack_notifier.Notify(user, repo, url, oldVer, newVer, GetUpdateLevel(oldVer, newVer), sendBool)
	}
}
