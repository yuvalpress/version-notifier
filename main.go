// version-notifier main
package main

import (
	"log"
	"os"
	anc "sirrend/internal/anchor"
	"sirrend/internal/commons"
	"sirrend/internal/config"
	"sirrend/internal/s3_client"
	smc "sirrend/internal/secrets_manager"
	"sirrend/internal/utils"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v2"
)

var (
	anchor anc.Anchor

	yamlData []byte

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

	// Get all secrets from Secret Manager
	log.Println("INFO: fetching secrets from AWS secret manager store.")
	versionNotifierSecret, exists := os.LookupEnv("SECRET_NAME_NOTIFIER")
	if !exists {
		log.Println("INFO: Could not file SECRET_NAME_NOTIFIER as env var.")
		os.Exit(1)
	}
	err := smc.ImportSecretsToEnv(versionNotifierSecret)
	if err != nil {
		log.Println("ERROR: Failed to get essential secret from AWS secrets store such as GITHUB_TOKEN or SLACK_TOKEN.")
		os.Exit(1)
	}

	// Get the config file from the S3 bucket
	svc, err := s3_client.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	yamlData = svc.GetObject(commons.NOTIFIER_BUCKET_PATH + commons.CONFIG_FILE_NAME)
	anchor.Init(yamlData)

	// Setup the notifications level for semver
	levels := utils.LevelsToNotify()
	log.Printf("INFO: Notifications will be sent for: %s\n", levels)

	// Set Log Level
	if LogLevel == "" {
		LogLevel = "INFO"
	}	
	log.Printf("INFO: Log Level is set for: %s\n", LogLevel)

	// Setup the Anchor Lists
	if len(anchor.RepoList) != 0 {
		log.Println("INFO: Core repository versions:")
	}
	for _, repoData := range anchor.RepoList {
		log.Printf("    %v/%v: %v\n", repoData.Owner, repoData.Project, repoData.Current)
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

func switchCurrentVersion(data []byte, owner string, project string, newVersion string) ([]byte, error) {
	// Unmarshal YAML data into the struct
	var newConf config.Conf
	err := yaml.Unmarshal(data, &newConf)
	if err != nil {
		return nil, err
	}

	//Iterate through the repos and update the currentVersion for the specified repoName
	for _, repo := range newConf.Repos {
		for key, value := range repo {
			if key == owner && value.Name == project {
				// Update the CurrentFlag field in the original struct
				value.CurrentFlag = newVersion
				// Update the map with the modified struct
				repo[key] = value
			}
		}
	}
	// Marshal the struct back to YAML
	newYamlData, err := yaml.Marshal(&newConf)

	if err != nil {
		return nil, err
	}
	if LogLevel == "DEBUG" {
		log.Printf("DEBUG: The new config file to be updated in S3: \n%s", newYamlData)
	}
	return newYamlData, nil
}

func updateConfigFile (newYamlData []byte) error {
	svc, err := s3_client.New()
	if err != nil {
		return nil
	}
	
	err = svc.UpdateObject(newYamlData, commons.NOTIFIER_BUCKET_PATH + commons.CONFIG_FILE_NAME)
	if err != nil {
		return err
	}

	log.Printf("INFO: Replaced new config file in the bucket: " + commons.NOTIFIER_BUCKET + "/" + commons.NOTIFIER_BUCKET_PATH + commons.CONFIG_FILE_NAME)
	return nil
}


// main - where the magic happens
func HandleRequest() {
	// initialize application
	levels := serviceInit()

	for index, repoData := range anchor.RepoList {
		latest, requestType, err := utils.GetVersion(repoData.Owner, repoData.Project)
		if err != nil {
			log.Printf(Red+" ERROR: Failed scraping %v: %v"+Reset, repoData.Owner+"/"+repoData.Project, err)
		}

		if latest != nil {
			var newVersion string
			if requestType == "release" {
				newVersion = utils.GetLatestTag(latest.Path("tag_name").String(), LogLevel)
			} else {
				newVersion = utils.GetLatestTag(latest.Path("name").String(), LogLevel)
			}

			result, newVer := utils.DoesNewTagExist(repoData.Current, newVersion, repoData.Owner+"/"+repoData.Project)

			if result {
				updateLevel := utils.GetUpdateLevel(repoData.Current, newVer)

				if utils.StringInSlice(updateLevel, levels) {
					if LogLevel == "DEBUG" || LogLevel == "INFO" {
						log.Printf(Green+"New %v version found for package %v/%v: %v\n"+Reset,
							updateLevel, repoData.Owner, repoData.Project, newVer)
					}

					// update releases link and create new config file
					var newURL string
					if requestType == "release" {
						newURL = strings.ReplaceAll(latest.Path("html_url").String(), "\"", "")
					} else {
						newURL = strings.ReplaceAll(latest.Path("zipball_url").String(), "\"", "")
					}

					//=============================================================//
					//=============================================================//
					newYamlData, err := switchCurrentVersion(yamlData, repoData.Owner, repoData.Project, newVersion)
					if err != nil {
						log.Println(err)
						os.Exit(3)
					}
					err = updateConfigFile(newYamlData)
					if err != nil {
						log.Println(err)
						os.Exit(3)
					}
					//=============================================================//
					// Updating the URL
					anchor.RepoList[index].URL = newURL

					// notify slack_notifier channel
					log.Println("INFO: Starting notifications services...")
					utils.Notify(repoData.Owner, repoData.Project, anchor.RepoList[index].URL, repoData.Current, "v"+newVer, requestType)
				}

				// update latest version
				anchor.RepoList[index].Current = "v" + newVer

			} else {
				log.Printf("INFO: No new version found for package %v/%v. Looks like already flagged to latest version: %v", repoData.Owner, repoData.Project, repoData.Current)
			}
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
