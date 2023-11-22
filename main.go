// version-notifier main
package main

import (
	"log"
	"os"
	anc "sirrend/internal/anchor"
	smc "sirrend/internal/secrets_manager"
	"sirrend/internal/utils"
	"strings"
)

var (
	anchor anc.Anchor

	// Reset color variables after call
	Reset = "\033[0m"

	// Green color for logs
	Green = "\033[32m"

	// Red color for logsss
	Red = "\033[31m"

	LogLevel = os.Getenv("LOG_LEVEL")
)

// serviceInit initializes the version-notifier service and returns the update levels to notify the client about
func serviceInit() []string {
	// initialize application data until successful
	log.Println("INFO: Starting application...")

	log.Println("INFO: Initializing latest tags for configured repositories")
	anchor.Init()

	levels := utils.LevelsToNotify()
	log.Printf("INFO: Notifications will be sent for: %s\n", levels)

	if LogLevel == "" {
		LogLevel = "INFO"
	}

	// Get all secrets from Secret Manager
	log.Println("INFO: fetching secrets from AWS secret manager store.")
	versionNotifierSecret, exists := os.LookupEnv("SECRET_NAME_TEST")
	if !exists {
		versionNotifierSecret = "SECRET"
	}

	err := smc.ImportSecretsToEnv(versionNotifierSecret)
	if err != nil {
		log.Println("ERROR: Failed to get essential secret from AWS secrets store such as GITHUB_TOKEN or SLACK_TOKEN.")
		os.Exit(1)
	}

	log.Printf("INFO: Log Level is set for: %s\n", LogLevel)

	if len(anchor.RepoList) != 0 {
		log.Println("INFO: Core repository versions:")
	}
	for _, repoData := range anchor.RepoList {
		log.Printf("    %v/%v: %v\n", repoData.User, repoData.Repo, repoData.Latest)
	}

	log.Println("INFO: Done!")
	log.Println("INFO: -----------------------------------------------------")

	// check if there are repos to scrape
	if len(anchor.RepoList) == 0 {
		log.Println(Red + "INFO: No repos to scrape... exiting" + Reset)
		os.Exit(1)
	}

	return levels
}

// main - where the magic happens
func main() {
	// initialize application
	levels := serviceInit()

	for index, repoData := range anchor.RepoList {
		latest, requestType, err := utils.GetVersion(repoData.User, repoData.Repo)
		if err != nil {
			log.Printf(Red+" ERROR: Failed scraping %v: %v"+Reset, repoData.User+"/"+repoData.Repo, err)
		}

		if latest != nil {
			log.Println("here")
			var newVersion string
			if requestType == "release" {
				newVersion = utils.GetLatestTag(latest.Path("tag_name").String(), LogLevel)
			} else {
				newVersion = utils.GetLatestTag(latest.Path("name").String(), LogLevel)
			}

			result, newVer := utils.DoesNewTagExist(repoData.Latest, newVersion, repoData.User+"/"+repoData.Repo)

			if result {
				updateLevel := utils.GetUpdateLevel(repoData.Latest, newVer)

				if utils.StringInSlice(updateLevel, levels) {
					if LogLevel == "DEBUG" || LogLevel == "INFO" {
						log.Printf(Green+"New %v version found for package %v/%v: %v\n"+Reset,
							updateLevel, repoData.User, repoData.Repo, newVer)
					}

					// update releases link
					var newURL string
					if requestType == "release" {
						newURL = strings.ReplaceAll(latest.Path("html_url").String(), "\"", "")
					} else {
						newURL = strings.ReplaceAll(latest.Path("zipball_url").String(), "\"", "")
					}

					anchor.RepoList[index].URL = newURL

					// notify slack_notifier channel
					utils.Notify(repoData.User, repoData.Repo, anchor.RepoList[index].URL, repoData.Latest, "v"+newVer, requestType)

				}

				// update latest version
				anchor.RepoList[index].Latest = "v" + newVer

			} else {
				if LogLevel == "DEBUG" {
					log.Printf("DEBUG: No new version found for package %v/%v", repoData.User, repoData.Repo)
				}
			}
		}
	}
}