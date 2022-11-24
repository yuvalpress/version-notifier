package utils

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	validate "golang.org/x/mod/semver"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Reset color variables after call
	Reset = "\033[0m"

	// Green color for logs
	Green = "\033[32m"

	// Red color for logs
	Red = "\033[31m"
)

// GetURL build the GitHub url with the needed user and repo
func GetURL(username, repoName string) string {
	return "https://api.github.com/repos/" + username + "/" + repoName + "/releases/latest"
}

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
		if value == level {
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
		log.Println(Red + "Wrong INTERVAL environment variable inserted, defaulting to 20 min" + Reset)
		return "30", 30
	}

	return interval, intInterval
}

func WaitForInterval() {
	_, intInterval := GetInterval()
	//set interval to run every <interval>
	for i := intInterval; i > 0; i-- {
		for s := 1; s <= 60; s++ {
			fmt.Printf("\r%v Performing next request in: %d minutes", time.Now().Format("2006/1/2 15:04:05"), i)
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Print("\n")
}
