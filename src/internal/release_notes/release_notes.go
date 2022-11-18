package release_notes

import (
	"github.com/go-rod/rod"
	"regexp"
	"strings"
)

// GetReleaseNotes receives url as string and fetches the release notes from it - if exist
func GetReleaseNotes(url string) (string, string) {
	blockType := "mrkdwn"
	page := rod.New().MustConnect().MustPage(url).MustWaitLoad()

	// fetch release notes div if exists
	check := page.MustHas("#repo-content-pjax-container > div > div > div > div.Box-body > div.markdown-body.my-3")
	if check {
		markdown := page.MustElement("#repo-content-pjax-container > div > div > div > div.Box-body > div.markdown-body.my-3").MustHTML()

		// find all tags with regex
		compile, _ := regexp.Compile("<.*?>")
		tags := compile.FindAllString(markdown, -1)

		if len(tags) > 0 {
			markdown = strings.Replace(markdown, tags[0], "", 1)
			markdown = strings.Replace(markdown, tags[len(tags)-1], "", 1)
		} else {
			blockType = "plain_text"
		}

		return blockType, markdown
	}

	return "", ""
}
