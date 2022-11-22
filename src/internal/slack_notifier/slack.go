package slack_notifier

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
	"yuvalpress/version-notifier/internal/release_notes"
)

var (
	// Reset color variables after call
	Reset = "\033[0m"

	// Red color for logs
	Red = "\033[31m"
)

// Notify sends a slack message with the supplied data
func Notify(user, repo, url, oldVer, newVer, updateLevel string, sendFullChangelog bool) {
	slackToken, exists := os.LookupEnv("SLACK_TOKEN")
	if !exists {
		log.Panicln("The SLACK_TOKEN environment variable doesn't exist")
	}

	slackChannel, exists := os.LookupEnv("SLACK_CHANNEL")
	if !exists {
		log.Panicln("The SLACK_CHANNEL environment variable doesn't exist")
	}

	slackClient := slack.New(slackToken)

	attachment := slack.Attachment{
		Pretext: "New Version Details:",
		Text:    url,
	}

	notes := release_notes.GetReleaseNotes(url, "text")
	if notes != "" && sendFullChangelog {
		_, _, err := slackClient.PostMessage(
			slackChannel,
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionBlocks(
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*New "+updateLevel+" update found for package: "+user+"/"+repo+"*"+"\n"+oldVer+" -> "+newVer, false, false), nil, nil),
				slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", notes, false, false), nil, nil)))

		if err != nil {
			fmt.Printf(Red+"Faild to post message to slack_notifier with the following error: %s\n"+Reset, err)
			return
		}

	} else {
		_, _, err := slackClient.PostMessage(
			slackChannel,
			slack.MsgOptionText("*New "+updateLevel+" update found for package: "+user+"/"+repo+"*"+"\n"+oldVer+" -> "+newVer, false),
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionUsername("Version Notifier"),
		)
		if err != nil {
			fmt.Printf(Red+"Faild to post message to slack_notifier with the following error: %s\n"+Reset, err)
			return
		}
	}
}
