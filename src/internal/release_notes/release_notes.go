package release_notes

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"strings"
)

// GetReleaseNotes receives url as string and fetches the release notes from it - if exist
func GetReleaseNotes(url string) string {
	divXPath := ""
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage(url).MustWaitLoad()

	// fetch release notes div if exists
	elements, _ := page.Elements("div")
	for _, v := range elements {
		if v.MustAttribute("Class") != nil {
			if strings.Contains(*v.MustAttribute("Class"), "markdown-body my-3") {
				divXPath = v.MustGetXPath(true)
			}
		}
	}

	if divXPath != "" {
		return page.MustElementX(divXPath).MustText()
	}

	return ""
}
