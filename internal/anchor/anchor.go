package anchor

import (
	"log"
	"os"
	"sirrend/internal/config"
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
	RepoList []Current
}

// Latest holds all the needed information for a repo instance
type Current struct {
	Owner   string
	Project string
	Current string
	URL     string
}

// Init method for main Anchor object
func (a *Anchor) Init(yamlData []byte) {
	confData, err := config.ReadConfigFile(yamlData)
	if err != nil {
		log.Fatalf("FATAL: Failed during initialization process with the following error: %v", err)
	}

	for _, repo := range confData.Repos {
		for ownerName, repoValues := range repo {
			project := repoValues.Name
			version := repoValues.CurrentFlag
			log.Println("INFO: Iterating over the " + ownerName + "/" + project + ":" + version + " project")
			latest := Current{ownerName, project, version, ""}
			log.Printf("INFO: Current state for the project is: %s" , latest)
			a.RepoList = append(a.RepoList, latest)
		}
	}
}