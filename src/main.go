// version-notifier main
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"yuvalpress/version-notifier/internal/scraper"
	"yuvalpress/version-notifier/internal/slack_notifier"
	"yuvalpress/version-notifier/internal/telegram_notifier"
	"yuvalpress/version-notifier/internal/utils"

	jparser "github.com/Jeffail/gabs/v2"
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

	LogLevel = os.Getenv("LOG_LEVEL")
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
			data, err := getVersion(username, repoName)
			if err != nil {
				log.Printf("Failed getting latest release of "+username+"/"+repoName+" with the following error: "+Red+"%v"+Reset, err)
				log.Println("Skipping..")
				continue
			}

			if data.Path("tag_name").String() == "" {
				return false
			}

			log.Println("Fetched latest release of: " + username + "/" + repoName)

			a.repoList = append(a.repoList, Latest{
				User:   username,
				Repo:   repoName,
				Latest: utils.GetLatestTag(data.Path("tag_name").String(), LogLevel),
				URL:    strings.ReplaceAll(data.Path("html_url").String(), "\"", ""),
			})
		}
	}

	return true
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

// get is responsible to fetch the latest data from the relative url
func getVersion(username, repoName string) (*jparser.Container, error) {
	url := utils.GetURL(username, repoName)
	request, err := scraper.APIRequest(url, LogLevel)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// notify is responsible for notifying a selected Slack channel.
// in the future, more methods will be added
func notify(user, repo, url, oldVer, newVer string) {
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
		telegram_notifier.Notify(user, repo, url, oldVer, newVer, utils.GetUpdateLevel(oldVer, newVer), sendBool)

	} else if method == "slack" {
		slack_notifier.Notify(user, repo, url, oldVer, newVer, utils.GetUpdateLevel(oldVer, newVer), sendBool)
	}
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

	levels := utils.LevelsToNotify()
	log.Printf("Notifications will be sent for: %s\n", levels)

	if LogLevel == "" {
		LogLevel = "INFO"
	}

	log.Printf("Log Level is set for: %s\n", LogLevel)

	interval, _ := utils.GetInterval()
	log.Printf("Interval is set to: %s minutes\n", interval)

	if len(anchor.repoList) != 0 {
		log.Println("Core repository versions:")
	}
	for _, repoData := range anchor.repoList {
		log.Printf("    %v/%v: %v\n", repoData.User, repoData.Repo, repoData.Latest)
	}

	log.Println("Done!")
	log.Println("-----------------------------------------------------")

	// check if there are repos to scrape
	if len(anchor.repoList) == 0 {
		log.Println(Red + "No repos to scrape... exiting" + Reset)
		os.Exit(1)
	}

	// loop to infinity
	for true {
		utils.WaitForInterval()
		for index, repoData := range anchor.repoList {
			latest, err := getVersion(repoData.User, repoData.Repo)
			if err != nil {
				log.Printf(Red+"Failed scraping %v: %v"+Reset, repoData.User+"/"+repoData.Repo, err)
			}

			if latest != nil {
				result, newVer := utils.DoesNewTagExist(repoData.Latest, utils.GetLatestTag(latest.Path("tag_name").String(), LogLevel), repoData.User+"/"+repoData.Repo)

				if result {
					updateLevel := utils.GetUpdateLevel(repoData.Latest, newVer)
					if utils.StringInSlice(updateLevel, levels) {
						if LogLevel == "DEBUG" || LogLevel == "INFO" {
							log.Printf(Green+"New %v version found for package %v/%v: %v\n"+Reset,
								updateLevel, repoData.User, repoData.Repo, newVer)
						}

						// update releases link
						anchor.repoList[index].URL = strings.ReplaceAll(latest.Path("html_url").String(), "\"", "")

						// notify slack_notifier channel
						notify(repoData.User, repoData.Repo, anchor.repoList[index].URL, repoData.Latest, "v"+newVer)

					}

					// update latest version
					anchor.repoList[index].Latest = "v" + newVer

				} else {
					if LogLevel == "DEBUG" {
						log.Printf("No new version found for package %v/%v", repoData.User, repoData.Repo)
					}
				}
			}
		}
	}

}
