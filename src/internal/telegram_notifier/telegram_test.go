package telegram_notifier

import (
	"os"
	"testing"
)

func TestNotify(t *testing.T) {
	// TODO: Fill this values using GitHub Actions
	_ = os.Setenv("TELEGRAM_TOKEN", "")
	_ = os.Setenv("TELEGRAM_CHAT_ID", "")

	Notify("google",
		"go-github",
		"https://github.com/google/go-github/releases/tag/v48.1.0",
		"v48.0.0", "v48.1.0", "minor", false)
}
