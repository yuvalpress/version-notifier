package telegram_notifier

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"sirrend/version-notifier/internal/release_notes"
	"strconv"
	"strings"
)

// getBot returns a telegram bot initialized with the token in context
func getBot() (*tg.BotAPI, error) {
	token, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		log.Panicln("The TELEGRAM_TOKEN environment variable doesn't exist")
	}

	bot, err := tg.NewBotAPI(token)
	bot.Debug = true

	if err != nil {
		log.Panicf("Failed creating telegram bot: %e", err)
	}

	return bot, nil
}

// sendMessage implements the basic use of sending a message and handles the error
func sendMessage(bot *tg.BotAPI, msg tg.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed sending message: %e\n", err)
	}
}

// sendReleaseNotes retrieves the release notes for the given release and sends them as text
func sendReleaseNotes(bot *tg.BotAPI, msg tg.MessageConfig, url string) {
	notes := release_notes.GetReleaseNotes(url, "text")
	if notes != "" {
		releaseNotesMessage := "\n*Release Notes:*\n" + strings.ReplaceAll(strings.ReplaceAll(notes, "_", "\\_"), "*", "\\*")

		msg.Text = releaseNotesMessage
		if len(msg.Text) > 4095 {
			for i := 0; i <= len(msg.Text); i = i + 4095 {
				msg.Text = msg.Text[i : i+4095]
				sendMessage(bot, msg)
			}
		} else {
			sendMessage(bot, msg)
		}
	}
}

// Notify notifies the telegram channel given in context
func Notify(user, repo, url, oldVer, newVer, updateLevel, versionType string, sendFullChangelog bool) {
	// configure bot and chatID
	bot, _ := getBot()
	chatID, exists := os.LookupEnv("TELEGRAM_CHAT_ID")
	if !exists {
		log.Panicln("The TELEGRAM_CHAT_ID environment variable doesn't exist")
	}

	IntChatID, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Panicln("The TELEGRAM_CHAT_ID environment variable cannot be converted to Int.\nPlease enter numbers only")
	}

	message := "*New " + updateLevel + " update found for package: " + user + "/" + repo + "*" + "\n" + oldVer + " -> " + newVer

	// initialize msg object
	msg := tg.NewMessage(IntChatID, message)
	msg.ParseMode = "Markdown"

	// send message
	if sendFullChangelog {
		sendMessage(bot, msg)
		if versionType == "release" {
			sendReleaseNotes(bot, msg, url)
		}

	} else {
		if versionType == "release" {
			msg.Text = msg.Text + "\n\n" + "*New Version Details:*\n" + url
		}
		sendMessage(bot, msg)
	}
}
