package anchor

import (
	jparser "github.com/Jeffail/gabs/v2"
	"log"
	"os"
	"strings"
	"yuvalpress/version-notifier/internal/config"
	"yuvalpress/version-notifier/internal/utils"
)

var (
	// Reset color variables after call
	Reset = "\033[0m"

	// Red color for logs
	Red = "\033[31m"

	LogLevel = os.Getenv("LOG_LEVEL")
)

// Anchor holds the first initialized information for the service
type Anchor struct {
	RepoList []Latest
}

// Latest holds all the needed information for a repo instance
type Latest struct {
	User   string
	Repo   string
	Latest string
	URL    string
}

func (l *Latest) init(t, username, repoName string, data *jparser.Container) {
	l.User = username
	l.Repo = repoName

	if t == "release" {
		l.Latest = utils.GetLatestTag(data.Path("tag_name").String(), LogLevel)
		l.URL = strings.ReplaceAll(data.Path("html_url").String(), "\"", "")
	} else if t == "tag" {
		l.Latest = utils.GetLatestTag(data.Path("name").String(), LogLevel)
		l.URL = strings.ReplaceAll(data.Path("zipball_url").String(), "\"", "")
	}

}

// Init method for main Anchor object
func (a *Anchor) Init() bool {
	confData, err := config.ReadConfigFile()
	if err != nil {
		log.Fatalf("Failed during initialization process with the following error: %v", err)
	}

	for _, info := range confData.Repos {
		for username, repoName := range info {
			data, requestType, err := utils.GetVersion(username, repoName)
			if err != nil {
				log.Printf("Failed getting latest release of "+username+"/"+repoName+" with the following error: "+Red+"%v"+Reset, err)
				log.Println("Skipping..")
				continue
			}

			if requestType == "release" && data.Path("tag_name").String() == "" || requestType == "tag" && data.Path("name").String() == "" {
				return false
			}

			log.Println("Fetched latest asset of: " + username + "/" + repoName)

			latest := Latest{}
			latest.init(requestType, username, repoName, data)
			a.RepoList = append(a.RepoList, latest)
		}
	}

	return true
}
