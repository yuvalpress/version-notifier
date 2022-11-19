package telegram_notifier

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"yuvalpress/version-notifier/internal/release_notes"
)

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

func Notify(user, repo, url, oldVer, newVer, updateLevel string) {
	bot, _ := getBot()
	notes := release_notes.GetReleaseNotes(url, "mrkdwn")
	chatID, exists := os.LookupEnv("TELEGRAM_CHAT_ID")
	if !exists {
		log.Panicln("The TELEGRAM_CHAT_ID environment variable doesn't exist")
	}

	IntChatID, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Panicln("The TELEGRAM_CHAT_ID environment variable cannot be converted to Int.\nPlease enter numbers only")
	}

	message := "*New " + updateLevel + " update found for package: " + user + "/" + repo + "*" + "\n" + oldVer + " -> " + newVer + "\n*ChangeLog:\n" + notes

	msg := tg.NewMessage(IntChatID, message)
	msg.ParseMode = "Markdown"

	_, err = bot.Send(msg)
	if err != nil {
		log.Panicf("Failed sending message: %e", err)
	}
}
