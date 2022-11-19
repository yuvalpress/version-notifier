package release_notes

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"regexp"
	"strings"
)

// GetReleaseNotes receives url as string and fetches the release notes from it - if exist
// format can be mrkdwn or text
func GetReleaseNotes(url, format string) string {
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
		if format == "text" {
			return page.MustElementX(divXPath).MustText()

		} else if format == "mrkdwn" {
			markdown := page.MustElementX(divXPath).MustHTML()
			// find all tags with regex
			compile, _ := regexp.Compile("<.*?>")
			tags := compile.FindAllString(markdown, -1)

			if len(tags) > 0 {
				markdown = strings.Replace(markdown, tags[0], "", 1)
				markdown = strings.Replace(markdown, tags[len(tags)-1], "", 1)
			}

			return markdown
		}
	}

	return ""
}
